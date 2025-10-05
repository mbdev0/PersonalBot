package trading

import (
	"context"
	"pump_fun/infrastructure/solana_price"
	"pump_fun/internal/core/position"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/monitoring"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/monitoring/models"
	positionservice "pump_fun/internal/services/position"
	positionhub "pump_fun/internal/services/subscription_hub/position"
	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/internal/services/trading/sell"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

type Strategy struct {
	taskService     *taskservice.TaskService
	positionHub     *positionhub.SubscriptionHub
	positionService *positionservice.Service
}

func (s *Strategy) NewTradingStrategy(ts *taskservice.TaskService, ph *positionhub.SubscriptionHub, ps *positionservice.Service) {
	s.taskService = ts
	s.positionHub = ph
	s.positionService = ps
}

func (s *Strategy) AfkSniping(afkTask *strategies.Afk, ctx context.Context) {
	coins := make(chan models.Coin, 100)
	defer close(coins)

	filterPipeline := filters.FilterPipeline{}
	for _, f := range afkTask.Filters {
		filterPipeline.AddFilter(f())
	}

	go monitoring.StartAFKMonitor(filterPipeline, coins, ctx)
	logger.Information("started afk monitor")

	for {
		select {
		case <-ctx.Done():
			return
		case coin, ok := <-coins:
			if !ok {
				return
			}
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
		s.monitorPositionForSellStrategies(*pos, afkTask, sellStrats, ctx)
	}

}
func (s *Strategy) createAndRunBuyTask(coin models.Coin, afkTask *strategies.Afk) (bt tasks.Task, err error) {
	coinAddr, err := solana.PublicKeyFromBase58(coin.CoinData.TokenAddr)
	if err != nil {
		logger.Error("couldn't read token address correctly: " + err.Error())
	}

	t := s.createBuyTask(afkTask, coinAddr)

	bt, err = s.taskService.Create(t)
	if err != nil {
		return nil, err
	}

	err = s.taskService.TransitionTask(bt.Id(), tasks.TaskRun)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

func (s *Strategy) createBuyTask(afkTask *strategies.Afk, tokenAddr solana.PublicKey) *tasks.BuyTask {
	bt := tasks.NewBuyTask(afkTask.Wallet, tokenAddr,
		[]tasks.Option{tasks.WithSlippage(afkTask.Slippage), tasks.WithComputeUnits(uint32(afkTask.ComputeUnits))},
		[]tasks.BuyOption{tasks.WithBuyAmount(afkTask.BuyAmount), tasks.WithBuyFee(afkTask.BuyFee)},
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

func (s *Strategy) monitorPositionForSellStrategies(pos position.Position, afkTask *strategies.Afk, strats []sell.Strategy, ctx context.Context) error {
	sub, err := s.positionHub.Subscribe(pos.PositionId, true)
	if err != nil {
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

		hasHit := s.handleSellStrategy(&msg, *afkTask, pos, strats, ctx)
		if hasHit {
			return nil
		}
	}
}

func (s *Strategy) handleSellStrategy(posMessage *position.PositionMessage, t strategies.Afk, pos position.Position, sellStrats []sell.Strategy, ctx context.Context) bool {
	for _, strategy := range sellStrats {
		hasHit := strategy.CheckIfPositionHasHit(posMessage)
		if hasHit && ctx.Err() == nil {
			s.createAndRunSellTask(t, &pos, strategy.SellAmount())
			return true
		}
	}

	return false
}

func (s *Strategy) createAndRunSellTask(afkTask strategies.Afk, pos *position.Position, sellAmount float64) {
	tsk, err := s.taskService.Create(tasks.NewSellTask(
		afkTask.Wallet,
		pos.TokenAddress,
		[]tasks.Option{
			tasks.WithComputeUnits(uint32(afkTask.ComputeUnits)),
			tasks.WithSlippage(afkTask.Slippage),
		},
		[]tasks.SellOption{
			tasks.WithSellAmount(sellAmount),
			tasks.WithSellFee(afkTask.SellFee),
			tasks.WithSellPositionId(pos.PositionId),
		},
	))

	if err != nil {
		logger.Error(err)
	}

	s.taskService.TransitionTask(tsk.Id(), tasks.TaskRun)
}
