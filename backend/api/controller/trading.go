package controller

import (
	"context"
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/api/mapper"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/services/subscription_hub/strategy"
	"personal_bot/internal/services/trading"
	"personal_bot/internal/services/wallet"
)

type StrategyController struct {
	strategyService *trading.Service
	walletService   *wallet.Service
}

func (sc *StrategyController) New(service *trading.Service, walletService *wallet.Service) {
	sc.strategyService = service
	sc.walletService = walletService
}

func (sc *StrategyController) Create(ctx context.Context, task dto.TradingTask) (*dto.TradingTaskResponse, error) {
	err := task.Validate()
	if err != nil {
		return nil, err
	}

	wallet, err := sc.walletService.GetByName(ctx, task.WalletName)
	if err != nil {
		return nil, err
	}

	t, err := mapper.MapTradingTaskDtoToTradingTask(task, wallet)
	if err != nil {
		return nil, err
	}

	createdTask, err := sc.strategyService.Create(t)
	if err != nil {
		return nil, err
	}

	resp, err := mapper.MapTradingTaskToDto(createdTask)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc *StrategyController) Delete(id int64, ctx context.Context) error {
	err := sc.strategyService.Delete(id, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sc *StrategyController) GetBy(id int64) (*dto.TradingTaskResponse, error) {
	task, err := sc.strategyService.GetBy(id)
	if err != nil {
		return nil, err
	}

	resp, err := mapper.MapTradingTaskToDto(task)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc *StrategyController) GetAll() ([]dto.TradingTaskResponse, error) {
	allTasks := sc.strategyService.GetAll()
	responseAllTasks := make([]dto.TradingTaskResponse, 0, len(allTasks))

	for _, val := range allTasks {
		resp, err := mapper.MapTradingTaskToDto(val)
		if err != nil {
			return nil, err
		}
		responseAllTasks = append(responseAllTasks, *resp)
	}

	return responseAllTasks, nil
}

func (sc *StrategyController) Update(ctx context.Context, id int64, tsk dto.TradingTaskPatch) (*dto.TradingTaskResponse, error) {
	task, err := sc.strategyService.GetBy(id)
	if err != nil {
		return nil, fmt.Errorf("task not found with the id %d", id)
	}

	//if the user has passed in a wallet into the patch - we will get it
	//this will then be passed into the mapper
	//if the wallet is nil (which it will be if the user didn't pass in the wallet in the patch)
	//we will not update the wallet inside the mapper
	var wallet *wallets.SolanaWallet
	if tsk.WalletName != nil {
		walletResp, err := sc.walletService.GetByName(ctx, *tsk.WalletName)
		if err != nil {
			return nil, err
		}

		wallet = &walletResp
	}

	patch, err := mapper.MapTradingTaskPatchDtoToTradingTaskPatch(tsk, dto.TradingType(task.StrategyType()), wallet)
	if err != nil {
		return nil, fmt.Errorf("error whilst mapping %w", err)
	}

	resp, err := sc.strategyService.Update(task, patch)
	if err != nil {
		return nil, err
	}

	mappedResp, err := mapper.MapTradingTaskToDto(resp)
	if err != nil {
		return nil, err
	}

	return mappedResp, nil
}

func (sc *StrategyController) Start(id int64) error {
	err := sc.strategyService.Start(id)
	if err != nil {
		return err
	}
	return nil
}

func (sc *StrategyController) Stop(id int64) error {
	err := sc.strategyService.Stop(id)
	if err != nil {
		return err
	}
	return nil
}

func (sc *StrategyController) Subscribe(id int64) (*strategy.Subscription, error) {
	tsk, err := sc.strategyService.GetBy(id)
	if err != nil {
		return nil, err
	}

	sub, err := sc.strategyService.Subscribe(tsk.StrategyTaskId())
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (sc *StrategyController) Unsubscribe(id int64) error {
	tsk, err := sc.strategyService.GetBy(id)
	if err != nil {
		return err
	}

	err = sc.strategyService.Unsubscribe(tsk.StrategyTaskId())
	if err != nil {
		return err
	}

	return nil

}
