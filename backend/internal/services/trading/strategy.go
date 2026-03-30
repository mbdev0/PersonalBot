package trading

import (
	"context"
	"fmt"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/position"
	rpcgroupsModel "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/monitoring"
	"personal_bot/internal/monitoring/filters"
	"personal_bot/internal/monitoring/models"
	positionservice "personal_bot/internal/services/position"
	rpcgroups "personal_bot/internal/services/rpc_groups"

	subscriptionhub "personal_bot/internal/services/subscription_hub"
	positionhub "personal_bot/internal/services/subscription_hub/position"
	"personal_bot/internal/services/subscription_hub/strategy"
	taskservice "personal_bot/internal/services/task_service"
	"personal_bot/internal/services/trading/sell"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Strategy struct {
	taskService     *taskservice.TaskService
	positionHub     *positionhub.SubscriptionHub
	strategyHub     *strategy.SubscriptionHub
	taskHub         *subscriptionhub.Hub
	positionService *positionservice.Service
	rpcService      *rpcgroups.Service
}

func NewTradingStrategy(ts *taskservice.TaskService, ph *positionhub.SubscriptionHub, ps *positionservice.Service, sh *strategy.SubscriptionHub, th *subscriptionhub.Hub, rs *rpcgroups.Service) *Strategy {
	return &Strategy{
		taskService:     ts,
		positionHub:     ph,
		positionService: ps,
		strategyHub:     sh,
		taskHub:         th,
		rpcService:      rs,
	}
}

func (s *Strategy) Sell(tsk *strategies.Sell, ctxCancel context.Context) {
	if err := s.taskService.StartTask(tsk.SellTaskId); err != nil {
		logger.Error(err)
	}

	go s.syncStateAndMessage(tsk.SellTaskId, tsk, ctxCancel)

}

func (s *Strategy) Buy(buyTask *strategies.Buy, ctx context.Context) {
	err := s.taskService.StartTask(buyTask.BuyTaskId)
	if err != nil {
		//TODO: when we're doing the PR for cleaning up WS messages - should set the task error here
		logger.Error(err)
	}

	go s.syncStateAndMessage(buyTask.BuyTaskId, buyTask, ctx)

	if len(buyTask.SellStrategies) != 0 {
		pos, ok := s.positionHub.WaitForCreate(buyTask.PositionId)
		if !ok {
			logger.Error("timeout whilst waiting for position to be created: ", buyTask.PositionId)
			return
		}

		sellStrats := s.ResolveStrategyConfig(buyTask.SellStrategies, pos, ctx)

		node, err := s.rpcService.GetNode(buyTask.RPCGroupId())
		if err != nil {
			logger.Error(err)
			return
		}

		err = s.monitorPositionForSellStrategies(*pos, buyTask, sellStrats, ctx, node)
		if err != nil {
			return
		}
	}

}

func (s *Strategy) syncStateAndMessage(taskId int64, strategyTask strategies.Task, ctx context.Context) error {
	task, err := s.taskService.GetTaskWith(taskId)
	if err != nil {
		logger.Error(err)
		return err
	}
	sub, err := s.taskHub.Subscribe(task)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Information("subscribed to task successfully")
	defer s.taskHub.Unsubcribe(task)

	for {
		select {
		case <-ctx.Done():
			logger.Information("strategy task stopped")
			strategyTask.SetStrategyState(string(strategies.CANCELLED))
			err := s.strategyHub.PublishStateUpdate(strategyTask.StrategyTaskId(), strategyTask.StrategyState())
			if err != nil {
				logger.Error(err)
				return err
			}
			return ctx.Err()
		case msg, ok := <-sub.Chan():
			if !ok {
				logger.Information("sub chan closed")
				return nil
			}
			s.processMessage(msg, strategyTask)
		}
	}
}

func (s *Strategy) processMessage(msg tasks.TaskEvent, strategyTask strategies.Task) error {
	if msg.EventType == tasks.StateUpdate {
		strategyTask.SetStrategyState(msg.State.TaskState)
		err := s.strategyHub.PublishStateUpdate(strategyTask.StrategyTaskId(), strategyTask.StrategyState())
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	if msg.EventType == tasks.ProgressMessage {
		strategyTask.SetStrategyMessage(msg.Message)
		err := s.strategyHub.PublishProgressMessage(strategyTask.StrategyTaskId(), strategyTask.StrategyMessage())
		if err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

func (s *Strategy) AfkSniping(afkTask *strategies.Afk, ctx context.Context) {
	coins := make(chan models.Coin, 100)
	defer close(coins)

	filterPipeline := filters.FilterPipeline{}
	for _, f := range afkTask.Filters {
		filterPipeline.AddFilter(f())
	}

	node, err := s.rpcService.GetNode(afkTask.RPCGroupId())
	if err != nil {
		logger.Error(err)
		return
	}

	go monitoring.StartAFKMonitor(filterPipeline, coins, ctx, node.WS)
	logger.Information("started afk monitor")

	for {
		select {
		case <-ctx.Done():
			afkTask.SetStrategyState(string(strategies.CANCELLED))
			err := s.strategyHub.PublishStateUpdate(afkTask.StrategyTaskId(), afkTask.StrategyState())
			if err != nil {
				logger.Error(err)
			}

			return
		case coin, ok := <-coins:
			if !ok {
				afkTask.SetStrategyState(string(strategies.CANCELLED))
				err := s.strategyHub.PublishStateUpdate(afkTask.StrategyTaskId(), afkTask.StrategyState())
				if err != nil {
					logger.Error(err)
				}
				return
			}
			logger.Information("found new coin: ", coin.CoinData.Name)
			s.handleNewCoin(afkTask, ctx, coin)
		}
	}
}

func (s *Strategy) handleNewCoin(afkTask *strategies.Afk, ctx context.Context, coin models.Coin) {
	bt, err := s.createAndRunBuyTask(coin, afkTask)
	if err != nil {
		logger.Error(err)
		return
	}

	if len(afkTask.SellStrategies) != 0 {
		pos, ok := s.positionHub.WaitForCreate(bt.Id())
		if !ok {
			logger.Error("timeout whilst waiting for position to be created: ", bt.Id())
			return
		}

		sellStrats := s.ResolveStrategyConfig(afkTask.SellStrategies, pos, ctx)
		node, err := s.rpcService.GetNode(afkTask.RPCGroupId())
		if err != nil {
			logger.Error(err)
			return
		}

		err = s.monitorPositionForSellStrategies(*pos, afkTask, sellStrats, ctx, node)
		if err != nil {
			return
		}
	}

}

func (s *Strategy) createAndRunBuyTask(coin models.Coin, afkTask *strategies.Afk) (bt tasks.Task, err error) {

	coinAddr, err := solana.PublicKeyFromBase58(coin.CoinData.TokenAddr)
	if err != nil {
		logger.Error("couldn't read token address correctly: " + err.Error())
	}

	node, err := s.rpcService.GetNode(afkTask.RPCGroupId())

	t := s.createBuyTask(afkTask, coinAddr, node)

	bt, err = s.taskService.Create(t)
	if err != nil {
		return nil, err
	}

	s.strategyHub.PublishTakeCreation(afkTask.StrategyTaskId(), bt)
	s.strategyHub.PublishProgressMessage(afkTask.StrategyTaskId(), fmt.Sprintf("created + running coin for %s", coin.CoinData.Symbol))

	err = s.taskService.StartTask(bt.Id())
	if err != nil {
		return nil, err
	}

	return bt, nil
}

func (s *Strategy) createBuyTask(afkTask *strategies.Afk, tokenAddr solana.PublicKey, rpc rpcgroupsModel.GroupItem) *tasks.BuyTask {
	bt := tasks.NewBuyTask(afkTask.Wallet, tokenAddr,
		[]tasks.Option{
			tasks.WithSlippage(afkTask.Slippage),
			tasks.WithComputeUnits(uint32(afkTask.ComputeUnits)),
			tasks.WithStrategyId(afkTask.StrategyTaskId()),
			tasks.WithHttpNode(rpc.Http),
			tasks.WithWS(rpc.WS),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(afkTask.BuyAmount),
			tasks.WithBuyFee(afkTask.BuyFee)},
	)
	return bt
}

func (s *Strategy) ResolveStrategyConfig(sellStratsConfig []strategies.StrategyConfig, position *position.Position, ctx context.Context) []sell.Strategy {
	strats := []sell.Strategy{}
	entryPrice, _ := position.EntryPrice.Float64()
	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		logger.Error(err)
	}

	//TODO: get total token amount as a call -> maybe from pos_sub
	totalTokenAmount := 1_000_000_000

	for _, strat := range sellStratsConfig {
		switch strat.Type {
		case strategies.StopLossMarketcap:
			newStrat := sell.NewStopLossMarketCap(strat.Value, *solPrice, float64(totalTokenAmount), strat.SellAmount)
			strats = append(strats, newStrat)
		case strategies.TakeProfitMarketCap:
			newStrat := sell.NewTakeProfitMarketCap(strat.Value, *solPrice, float64(totalTokenAmount), strat.SellAmount)
			strats = append(strats, newStrat)
		case strategies.StopLossPercentage:
			strats = append(strats, sell.NewStopLossPercentage(entryPrice, strat.Value, strat.SellAmount))
		case strategies.TakeProfitPercentage:
			strats = append(strats, sell.NewTakeProfitPercentage(entryPrice, strat.Value, strat.SellAmount))
		case strategies.StopLossPrice:
			strats = append(strats, sell.NewStopLossPrice(strat.Value, strat.SellAmount))
		case strategies.TakeProfitPrice:
			strats = append(strats, sell.NewTakeProfitPrice(strat.Value, strat.SellAmount))
		}
	}
	return strats
}

func (s *Strategy) monitorPositionForSellStrategies(pos position.Position, sellableTask strategies.SellableStrategy, strats []sell.Strategy, ctx context.Context, rpcNode rpcgroupsModel.GroupItem) error {
	rpc := rpcgroupsModel.RPCNode{
		Http: rpc.New(rpcNode.Http),
		WS:   rpcNode.WS,
	}
	sub, err := s.positionHub.Subscribe(pos.PositionId, true, &rpc)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		err := s.positionHub.Unsubscribe(pos.PositionId, true)
		if err != nil {
			logger.Error(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		msg, ok := <-sub.SubChan
		if !ok {
			return nil
		}

		if msg.MessageType == position.Stopped {
			return nil
		}

		hasHit := s.handleSellStrategy(&msg, sellableTask, pos, strats, ctx, rpcNode)
		if hasHit {
			return nil
		}
	}
}

func (s *Strategy) handleSellStrategy(posMessage *position.PositionMessage, sellableTask strategies.SellableStrategy, pos position.Position, sellStrats []sell.Strategy, ctx context.Context, rpcNode rpcgroupsModel.GroupItem) bool {
	for _, strategy := range sellStrats {
		hasHit := strategy.CheckIfPositionHasHit(posMessage)
		if hasHit && ctx.Err() == nil {
			s.createAndRunSellTask(sellableTask, &pos, strategy.SellAmount(), rpcNode)
			return true
		}
	}

	return false
}

func (s *Strategy) createAndRunSellTask(sellableTask strategies.SellableStrategy, pos *position.Position, sellAmount float64, rpcNode rpcgroupsModel.GroupItem) {
	tsk, err := s.taskService.Create(tasks.NewSellTask(
		sellableTask.GetWallet(),
		pos.TokenAddress,
		[]tasks.Option{
			tasks.WithComputeUnits(uint32(sellableTask.GetComputeUnits())),
			tasks.WithSlippage(sellableTask.GetSlippage()),
			tasks.WithStrategyId(sellableTask.StrategyTaskId()),
			tasks.WithHttpNode(rpcNode.Http),
			tasks.WithWS(rpcNode.WS),
		},
		[]tasks.SellOption{
			tasks.WithSellAmount(sellAmount),
			tasks.WithSellFee(*sellableTask.GetSellFee()),
			tasks.WithSellPositionId(&pos.PositionId),
		},
	))

	if err != nil {
		logger.Error(err)
	}
	s.taskService.StartTask(tsk.Id())
}
