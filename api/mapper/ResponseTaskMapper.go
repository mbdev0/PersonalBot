package mapper

import (
	"fmt"
	"pump_fun/api/models"
	"pump_fun/internal/models/tasks"
)

func MapTaskToReponseTask(task tasks.Task) (*models.ResponseTask, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return mapBuyToResponseTask(t), nil
	case *tasks.SellTask:
		return mapSellToResponseTask(t), nil
	default:
		return nil, fmt.Errorf("no mapper created for type: " + t.GetTaskType())
	}
}

func mapBuyToResponseTask(t *tasks.BuyTask) *models.ResponseTask {
	responseTask := models.ResponseTask{}
	responseTask.Type = t.TaskType
	responseTask.ComputeUnits = t.ComputeUnits
	responseTask.Slippage = t.Slippage
	responseTask.State.TaskState = t.State.TaskState.ToString()
	responseTask.State.Error = t.State.Error
	responseTask.TokenAddress = t.TokenAddress.String()
	responseTask.TaskId = t.TaskId
	responseTask.BuyFee = &t.BuyFee
	buyAmount, _ := t.BuyAmount.Float64()
	responseTask.BuyAmount = &buyAmount

	return &responseTask
}

func mapSellToResponseTask(t *tasks.SellTask) *models.ResponseTask {
	responseTask := models.ResponseTask{}
	responseTask.Type = t.TaskType
	responseTask.ComputeUnits = t.ComputeUnits
	responseTask.Slippage = t.Slippage
	responseTask.State.TaskState = t.State.TaskState.ToString()
	responseTask.State.Error = t.State.Error
	responseTask.TokenAddress = t.TokenAddress.String()
	responseTask.TaskId = t.TaskId
	responseTask.SellAmount = &t.PercentageToSell
	responseTask.SellFee = &t.SellFee

	return &responseTask
}
