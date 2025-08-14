package mapper

import (
	"fmt"
	"pump_fun/api/dto"
	"pump_fun/internal/core/tasks"
)

func MapTaskToReponseTask(task tasks.Task) (*dto.ResponseTask, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return mapBuyToResponseTask(t), nil
	case *tasks.SellTask:
		return mapSellToResponseTask(t), nil
	default:
		return nil, fmt.Errorf("no mapper created for type: " + t.Type())
	}
}

func mapBuyToResponseTask(t *tasks.BuyTask) *dto.ResponseTask {
	responseTask := dto.ResponseTask{}
	responseTask.Type = t.Type()
	responseTask.ComputeUnits = t.ComputeUnits()
	responseTask.Slippage = t.Slippage()
	state := t.State()
	responseTask.State.TaskState = state.TaskState.ToString()
	responseTask.State.Error = state.Error
	responseTask.TokenAddress = t.Token().String()
	responseTask.TaskId = t.Id()
	fee := t.Fee()
	responseTask.BuyFee = &fee
	buyAmount := t.BuyAmount()
	buyAmountFloat, _ := buyAmount.Float64()
	responseTask.BuyAmount = &buyAmountFloat

	return &responseTask
}

func mapSellToResponseTask(t *tasks.SellTask) *dto.ResponseTask {
	responseTask := dto.ResponseTask{}
	responseTask.Type = t.Type()
	responseTask.ComputeUnits = t.ComputeUnits()
	responseTask.Slippage = t.Slippage()
	state := t.State()
	responseTask.State.TaskState = state.TaskState.ToString()
	responseTask.State.Error = state.Error
	responseTask.TokenAddress = t.Token().String()
	responseTask.TaskId = t.Id()
	sellAmnt := t.SellPercentage()
	fee := t.Fee()
	responseTask.SellAmount = &sellAmnt
	responseTask.SellFee = &fee

	return &responseTask
}
