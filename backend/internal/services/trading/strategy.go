package trading

import (
	"context"
	"fmt"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/program"
	rpcgroupsModel "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/tasks"
	positionservice "personal_bot/internal/services/position"
	rpcgroups "personal_bot/internal/services/rpc_groups"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/monitoring/models"

	positionhub "personal_bot/internal/services/subscription_hub/position"
	"personal_bot/internal/services/subscription_hub/strategy"
	"personal_bot/internal/services/subscription_hub/task"
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
	taskHub         *task.Hub
	positionService *positionservice.Service
	rpcService      *rpcgroups.Service
}

func NewTradingStrategy(ts *taskservice.TaskService, ph *positionhub.SubscriptionHub, ps *positionservice.Service,
	sh *strategy.SubscriptionHub, th *task.Hub, rs *rpcgroups.Service) *Strategy {
	return &Strategy{
		taskService:     ts,
		positionHub:     ph,
		positionService: ps,
		strategyHub:     sh,
		taskHub:         th,
		rpcService:      rs,
	}
}

func (s *Strategy) Sell(ctxCancel context.Context, tsk *strategies.Sell) {
	if err := s.taskService.StartTask(tsk.SellTaskId); err != nil {
		logger.Error(err)
	}

	go s.syncStateAndMessage(ctxCancel, tsk.SellTaskId, tsk)

}

func (s *Strategy) Buy(ctx context.Context, buyTask *strategies.Buy) {
	err := s.taskService.StartTask(buyTask.BuyTaskId)
	if err != nil {
		buyTask.SetStrategyMessage(err.Error())
		logger.Error(err)
	}

	go s.syncStateAndMessage(ctx, buyTask.BuyTaskId, buyTask)

	if len(buyTask.SellStrategies) != 0 {
		pos, ok := s.positionHub.WaitForCreate(buyTask.PositionId)
		if !ok {
			logger.Error("timeout whilst waiting for position to be created: ", buyTask.PositionId)
			return
		}

		sellStrats := s.ResolveStrategyConfig(ctx, buyTask.SellStrategies, pos)

		node, err := s.rpcService.GetNode(buyTask.RPCGroupId())
		if err != nil {
			logger.Error(err)
			return
		}

		err = s.monitorPositionForSellStrategies(ctx, *pos, buyTask, sellStrats, node)
		if err != nil {
			return
		}
	}

}

func (s *Strategy) syncStateAndMessage(ctx context.Context, taskId int64, strategyTask strategies.Task) error {
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
	defer s.taskHub.Unsubcribe(task)

	for {
		select {
		case <-ctx.Done():
			strategyTask.SetStrategyState(string(strategies.CANCELLED))
			s.strategyHub.PublishStateUpdate(strategyTask.StrategyTaskId(), strategyTask.StrategyState())
			return ctx.Err()
		case msg, ok := <-sub.Chan():
			if !ok {
				return nil
			}
			s.processMessage(msg, strategyTask)
		}
	}
}

func (s *Strategy) processMessage(msg tasks.TaskEvent, strategyTask strategies.Task) {
	if msg.EventType == tasks.StateUpdate {
		strategyTask.SetStrategyState(msg.State.TaskState)
		s.strategyHub.PublishStateUpdate(strategyTask.StrategyTaskId(), strategyTask.StrategyState())
	}

	if msg.EventType == tasks.ProgressMessage {
		strategyTask.SetStrategyMessage(msg.Message)
		s.strategyHub.PublishProgressMessage(strategyTask.StrategyTaskId(), strategyTask.StrategyMessage())
	}
}

func (s *Strategy) AfkSniping(ctx context.Context, afkTask *strategies.Afk) {
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

	prog := program.Resolve(afkTask.Program)
	coinStream := prog.NewCoinStream(ctx, rpcgroupsModel.RPCNode{
		Http: rpc.New(node.Http),
		WS:   node.WS,
	}, filterPipeline)

	logger.Information("started afk monitor")

	counter := 0
	for {
		if counter > 3 {
			break
		}

		select {
		case <-ctx.Done():
			afkTask.SetStrategyState(string(strategies.CANCELLED))
			s.strategyHub.PublishStateUpdate(afkTask.StrategyTaskId(), afkTask.StrategyState())

			return
		case coin, ok := <-coinStream:
			if !ok {
				afkTask.SetStrategyState(string(strategies.CANCELLED))
				s.strategyHub.PublishStateUpdate(afkTask.StrategyTaskId(), afkTask.StrategyState())
				return
			}
			logger.Information("found new coin: ", coin.CoinData.Name)
			s.handleNewCoin(ctx, afkTask, coin)
			counter++
		}
	}
}

func (s *Strategy) handleNewCoin(ctx context.Context, afkTask *strategies.Afk, coin models.Coin) {
	bt, err := s.createAndRunBuyTask(ctx, coin, afkTask)
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

		sellStrats := s.ResolveStrategyConfig(ctx, afkTask.SellStrategies, pos)
		node, err := s.rpcService.GetNode(afkTask.RPCGroupId())
		if err != nil {
			logger.Error(err)
			return
		}

		err = s.monitorPositionForSellStrategies(ctx, *pos, afkTask, sellStrats, node)
		if err != nil {
			return
		}
	}

}

func (s *Strategy) createAndRunBuyTask(ctx context.Context, coin models.Coin, afkTask *strategies.Afk) (bt tasks.Task, err error) {

	coinAddr, err := solana.PublicKeyFromBase58(coin.CoinData.TokenAddr)
	if err != nil {
		logger.Error("couldn't read token address correctly: " + err.Error())
	}

	node, err := s.rpcService.GetNode(afkTask.RPCGroupId())
	if err != nil {
		return nil, err
	}

	_, err = s.rpcService.Load(ctx, afkTask.RPCGroupId())
	if err != nil {
		return nil, err
	}

	t := s.createBuyTask(afkTask, coinAddr, node)

	bt, err = s.taskService.Create(t)
	if err != nil {
		return nil, err
	}

	s.strategyHub.PublishTaskCreation(afkTask.StrategyTaskId(), bt)
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
			tasks.WithProgram(afkTask.GetProgram()),
			tasks.WithSlippage(afkTask.Slippage),
			tasks.WithComputeUnits(uint32(afkTask.ComputeUnits)),
			tasks.WithStrategyId(afkTask.StrategyTaskId()),
			tasks.WithRPCGroupId(afkTask.RPCGroupId()),
			tasks.WithHttpNode(rpc.Http),
			tasks.WithWS(rpc.WS),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(afkTask.BuyAmount),
			tasks.WithBuyFee(afkTask.BuyFee)},
	)
	return bt
}

func (s *Strategy) ResolveStrategyConfig(ctx context.Context, sellStratsConfig []strategies.StrategyConfig, position *position.Position) []sell.Strategy {
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

func (s *Strategy) monitorPositionForSellStrategies(ctx context.Context, pos position.Position, sellableTask strategies.SellableStrategy, strats []sell.Strategy, rpcNode rpcgroupsModel.GroupItem) error {
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
		s.positionHub.Unsubscribe(pos.PositionId, true)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-sub.SubChan:
			if !ok {
				return nil
			}

			if msg.MessageType == position.Stopped {
				return nil
			}

			hasHit := s.handleSellStrategy(ctx, &msg, sellableTask, pos, strats, rpcNode)
			if hasHit {
				return nil
			}
		}
	}
}

func (s *Strategy) handleSellStrategy(ctx context.Context, posMessage *position.PositionMessage, sellableTask strategies.SellableStrategy, pos position.Position, sellStrats []sell.Strategy, rpcNode rpcgroupsModel.GroupItem) bool {
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
			tasks.WithProgram(sellableTask.GetProgram()),
			tasks.WithComputeUnits(uint32(sellableTask.GetComputeUnits())),
			tasks.WithSlippage(sellableTask.GetSlippage()),
			tasks.WithStrategyId(sellableTask.StrategyTaskId()),
			tasks.WithRPCGroupId(sellableTask.RPCGroupId()),
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
