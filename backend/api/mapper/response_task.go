package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/tasks"
)

func MapTaskToReponseTask(task tasks.Task) (*dto.ResponseTask, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return mapBuyToResponseTask(t), nil
	case *tasks.SellTask:
		return mapSellToResponseTask(t), nil
	default:
		return nil, fmt.Errorf("no mapper created for type: %s", t.Type())
	}
}

func mapBuyToResponseTask(t *tasks.BuyTask) *dto.ResponseTask {
	responseTask := dto.ResponseTask{}
	responseTask.Type = t.Type()
	responseTask.ComputeUnits = t.ComputeUnits
	responseTask.Slippage = t.Slippage
	state := t.State()
	responseTask.State.TaskState = state.TaskState.ToString()
	responseTask.State.Error = state.Error
	responseTask.TokenAddress = t.Token.String()
	responseTask.TaskId = t.Id()
	fee := t.Fee
	responseTask.BuyFee = &fee
	buyAmount := t.BuyAmount
	buyAmountFloat, _ := buyAmount.Float64()
	convertedFloat := buyAmountFloat / constants.LamportsConversion
	responseTask.BuyAmount = &convertedFloat
	responseTask.WalletAddressName = t.WalletName
	responseTask.WalletPublicKey = t.Wallet.PublicKey().Short(constants.ShortPublicAddressInt)
	responseTask.StrategyId = t.StrategyId
	responseTask.TimeCreated = t.TimeCreated
	responseTask.Message = t.Message()
	responseTask.HTTPNode = t.HttpNode()
	responseTask.WSNode = t.WSNode()
	responseTask.RPCGroupId = t.GetRPCGroupId()

	return &responseTask
}

func mapSellToResponseTask(t *tasks.SellTask) *dto.ResponseTask {
	responseTask := dto.ResponseTask{}
	responseTask.Type = t.Type()
	responseTask.ComputeUnits = t.ComputeUnits
	responseTask.Slippage = t.Slippage
	state := t.State()
	responseTask.State.TaskState = state.TaskState.ToString()
	responseTask.State.Error = state.Error
	responseTask.TokenAddress = t.Token.String()
	responseTask.TaskId = t.Id()
	sellAmnt := t.SellPercentage
	fee := t.Fee
	responseTask.SellAmount = &sellAmnt
	responseTask.SellFee = &fee
	responseTask.SellPositionId = t.Position_id
	responseTask.WalletAddressName = t.WalletName
	responseTask.WalletPublicKey = t.Wallet.PublicKey().Short(constants.ShortPublicAddressInt)
	responseTask.TimeCreated = t.TimeCreated
	responseTask.StrategyId = t.StrategyId
	responseTask.Message = t.Message()
	responseTask.HTTPNode = t.HttpNode()
	responseTask.WSNode = t.WSNode()
	responseTask.RPCGroupId = t.GetRPCGroupId()

	return &responseTask
}
