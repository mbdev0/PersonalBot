package strategy

import (
	"context"
	"fmt"
	"personal_bot/backend/infrastructure/solana_price"
	"personal_bot/backend/internal/core/position"
	rpcgroupsModel "personal_bot/backend/internal/core/rpc_groups"
	"personal_bot/backend/internal/core/strategies"
	"personal_bot/backend/internal/core/tasks"
	positionservice "personal_bot/backend/internal/services/position"
	rpcgroups "personal_bot/backend/internal/services/rpc_groups"
	positionhub "personal_bot/backend/internal/services/subscription_hub/position"
	"personal_bot/backend/internal/services/subscription_hub/strategy"
	"personal_bot/backend/internal/services/subscription_hub/task"
	taskservice "personal_bot/backend/internal/services/task_service"
	"personal_bot/backend/internal/services/trading/sell"
	"personal_bot/backend/pkg/logger"

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

func New(ts *taskservice.TaskService, ph *positionhub.SubscriptionHub, ps *positionservice.Service,
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

func (s *Strategy) Run(ctx context.Context, strategyTask strategies.Task) error {
	strategyTask.SetStrategyState(string(strategies.RUNNING))
	s.strategyHub.PublishStateUpdate(strategyTask.StrategyTaskId(), strategyTask.StrategyState())

	switch tsk := strategyTask.(type) {
	case *strategies.Afk:
		afk := NewAFKEngine(*s)
		go afk.Run(ctx, tsk)
	case *strategies.Buy:
		buy := NewBuyEngine(*s)
		go buy.Run(ctx, tsk)
	case *strategies.Sell:
		sell := NewSellEngine(*s)
		go sell.Run(ctx, tsk)
	case *strategies.Spam:
		spam := NewSpamEngine(*s)
		go spam.Run(ctx, tsk)
	default:
		strategyTask.SetStrategyState(string(strategies.FAILED))
		s.strategyHub.PublishStateUpdate(strategyTask.StrategyTaskId(), strategyTask.StrategyState())
		return fmt.Errorf("task doesn't belong to a strategy")
	}

	return nil
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

func (s *Strategy) resolveStrategyConfig(sellStratsConfig []strategies.StrategyConfig, position *position.Position) []sell.Strategy {
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
