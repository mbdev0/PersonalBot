package controller

import (
	"fmt"
	"pump_fun/api/dto"
	"pump_fun/api/mapper"
	"pump_fun/internal/services/trading"
)

type StrategyController struct {
	strategyService *trading.Service
}

func (sc *StrategyController) New(service *trading.Service) {
	sc.strategyService = service
}

func (sc *StrategyController) Create(task dto.TradingTask) (*dto.TradingTaskResponse, error) {
	t, err := mapper.MapTradingTaskDtoToTradingTask(task)
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

func (sc *StrategyController) Delete(id string) error {
	err := sc.strategyService.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (sc *StrategyController) GetBy(id string) (*dto.TradingTaskResponse, error) {
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

func (sc *StrategyController) Update(id string, tsk dto.TradingTaskPatch) (*dto.TradingTaskResponse, error) {
	task, err := sc.strategyService.GetBy(id)
	if err != nil {
		return nil, fmt.Errorf("task not found with the id %s", id)
	}

	patch, err := mapper.MapTradingTaskPatchDtoToTradingTaskPatch(tsk, dto.TradingType(task.StrategyType()))
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

func (sc *StrategyController) Start(id string) {
}

func (sc *StrategyController) Stop(id string) {
}
