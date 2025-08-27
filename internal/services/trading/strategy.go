package trading

import (
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/monitoring"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/monitoring/models"
	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

type Service struct {
	taskService *taskservice.TaskService
}

func (s *Service) NewService(ts *taskservice.TaskService) {
	s.taskService = ts
}

func (s *Service) AfkSniping(afkTask strategies.Afk) {
	coins := make(chan models.Coin, 100)
	defer close(coins)

	// how do we pass filterPipeline to the AFK Monitor?
	// create filterPipeline pipeline here then pass it in?
	filterPipeline := filters.FilterPipeline{}
	for _, f := range afkTask.Filters {
		filterPipeline.AddFilter(f())
	}

	monitoring.StartAFKMonitor(filterPipeline, coins)

	for coin := range coins {
		coinAddr, err := solana.PublicKeyFromBase58(coin.CoinData.TokenAddr)
		if err != nil {
			logger.Error("couldn't read token address correctly: " + err.Error())
		}
		t := s.createBuyTask(strategies.Afk{}, coinAddr)

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

func (s *Service) createBuyTask(afkTask strategies.Afk, tokenAddr solana.PublicKey) *tasks.BuyTask {
	// buyTask := tasks.BuyTask{}
	bt := tasks.NewBuyTask(afkTask.Wallet, tokenAddr,
		[]tasks.Option{tasks.WithSlippage(afkTask.Slippage), tasks.WithComputeUnits(uint32(afkTask.ComputeUnits))},
		[]tasks.BuyOption{tasks.WithBuyAmount(afkTask.BuyAmount), tasks.WithBuyFee(afkTask.BuyFee)},
	)
	return bt
}
