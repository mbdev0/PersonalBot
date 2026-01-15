package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/solana/utils"
)

func MapTradingTaskPatchDtoToTradingTaskPatch(src dto.TradingTaskPatch, tskType dto.TradingType, wallet *wallets.SolanaWallet) (strategies.Patch, error) {

	var taskType dto.TradingType
	if src.Type != nil {
		taskType = *src.Type
	} else {
		taskType = dto.TradingType(tskType)
	}

	switch taskType {
	case dto.AFK:
		afkPatch, err := createAfkPatch(src, wallet)
		if err != nil {
			return nil, err
		}
		return afkPatch, nil
	default:
		return nil, fmt.Errorf("task type hasn't been set up for the type: %s", taskType)
	}

}

func createAfkPatch(src dto.TradingTaskPatch, wallet *wallets.SolanaWallet) (resp *strategies.AfkPatch, err error) {
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

	if src.SellFee != nil {
		respPatch.SellFee = src.SellFee
	}

	//if the wallet was set in the patch, we get it in the controller and pass it down
	//this is to stop mixing of services and keep them as top level as possible
	if wallet != nil {
		respPatch.Wallet = wallet
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
