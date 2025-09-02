package trading

import (
	"context"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/monitoring"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/monitoring/models"
	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

type Strategy struct {
	taskService *taskservice.TaskService
}

func (s *Strategy) NewTradingStrategy(ts *taskservice.TaskService) {
	s.taskService = ts
}

func (s *Strategy) AfkSniping(afkTask *strategies.Afk, ctx context.Context) {
	coins := make(chan models.Coin, 100)
	defer close(coins)

	filterPipeline := filters.FilterPipeline{}
	for _, f := range afkTask.Filters {
		filterPipeline.AddFilter(f())
	}

	go monitoring.StartAFKMonitor(filterPipeline, coins)
	logger.Information("started afk monitor")

	for coin := range coins {
		coinAddr, err := solana.PublicKeyFromBase58(coin.CoinData.TokenAddr)
		if err != nil {
			logger.Error("couldn't read token address correctly: " + err.Error())
		}
		t := s.createBuyTask(afkTask, coinAddr)

		bt, err := s.taskService.Create(t)
		if err != nil {
			logger.Error("error whilst creating task: ", err)
		}

		err = s.taskService.TransitionTask(bt.Id(), tasks.TaskRun)
		if err != nil {
			logger.Error("error whilst transitioning task: ", err)
			return
		}
	}
}

func (s *Strategy) createBuyTask(afkTask *strategies.Afk, tokenAddr solana.PublicKey) *tasks.BuyTask {
	// buyTask := tasks.BuyTask{}
	bt := tasks.NewBuyTask(afkTask.Wallet, tokenAddr,
		[]tasks.Option{tasks.WithSlippage(afkTask.Slippage), tasks.WithComputeUnits(uint32(afkTask.ComputeUnits))},
		[]tasks.BuyOption{tasks.WithBuyAmount(afkTask.BuyAmount), tasks.WithBuyFee(afkTask.BuyFee)},
	)
	return bt
}
