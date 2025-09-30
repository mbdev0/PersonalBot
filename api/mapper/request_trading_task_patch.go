package mapper

import (
	"fmt"
	"pump_fun/api/dto"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapTradingTaskPatchDtoToTradingTaskPatch(src dto.TradingTaskPatch, tskType dto.TradingType) (strategies.Patch, error) {

	var taskType dto.TradingType
	if src.Type != nil {
		taskType = *src.Type
	} else {
		taskType = dto.TradingType(tskType)
	}

	switch taskType {
	case dto.AFK:
		afkPatch, err := createAfkPatch(src)
		if err != nil {
			return nil, err
		}
		return afkPatch, nil
	default:
		return nil, fmt.Errorf("task type hasn't been set up for the type: %s", taskType)
	}

}

func createAfkPatch(src dto.TradingTaskPatch) (resp *strategies.AfkPatch, err error) {
	respPatch := strategies.AfkPatch{}
	if src.BuyAmount != nil {
		respPatch.BuyAmount = utils.ConvertSolToLamport(*src.BuyAmount)
	}

	if src.BuyFee != nil {
		respPatch.BuyFee = src.BuyFee
	}

	if src.ComputeUnits != nil {
		respPatch.ComputeUnits = src.ComputeUnits
	}

	if src.Slippage != nil {
		respPatch.Slippage = src.Slippage
	}

	if src.Wallet != nil {
		wallet, err := solana.PrivateKeyFromBase58(*src.Wallet)
		if err != nil {
			return nil, err
		}
		respPatch.Wallet = &wallet
	}

	if src.Filters != nil {
		filters := mapFiltersToDestFilters(*src.Filters)
		respPatch.Filters = &filters
	}

	if src.SellStrategies != nil {
		strats := mapDTOToStrategyConfigs(*src.SellStrategies)
		respPatch.SellStrategies = &strats
	}

	return &respPatch, nil
}
