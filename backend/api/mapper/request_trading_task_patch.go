package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/wallets"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapTradingTaskPatchDtoToTradingTaskPatch(src dto.TradingTaskPatch, tskType dto.TradingType, wallet *wallets.SolanaWallet, rpcGroup *rpcgroups.RPCGroup) (strategies.Patch, error) {

	var taskType dto.TradingType
	if src.Type != nil {
		taskType = *src.Type
	} else {
		taskType = dto.TradingType(tskType)
	}

	switch taskType {
	case dto.AFK:
		return createAfkPatch(src, wallet, rpcGroup)
	case dto.BUY:
		return createBuyPatch(src, wallet, rpcGroup)
	case dto.SELL:
		return createSellPatch(src, wallet, rpcGroup)
	default:
		return nil, fmt.Errorf("task type hasn't been set up for the type: %s", taskType)
	}

}

func createAfkPatch(src dto.TradingTaskPatch, wallet *wallets.SolanaWallet, rpcGroup *rpcgroups.RPCGroup) (resp *strategies.AfkPatch, err error) {
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

	respPatch.SellFee = src.SellFee

	//if the wallet was set in the patch, we get it in the controller and pass it down
	//this is to stop mixing of services and keep them as top level as possible
	if wallet != nil {
		respPatch.Wallet = wallet
	}

	if rpcGroup != nil {
		respPatch.RPCGroup = rpcGroup
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

func createBuyPatch(src dto.TradingTaskPatch, wallet *wallets.SolanaWallet, rpcGroup *rpcgroups.RPCGroup) (strategies.Patch, error) {
	respPatch := strategies.BuyPatch{}
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

	respPatch.SellFee = src.SellFee

	if rpcGroup != nil {
		respPatch.RPCGroup = rpcGroup
	}

	//if the wallet was set in the patch, we get it in the controller and pass it down
	//this is to stop mixing of services and keep them as top level as possible
	if wallet != nil {
		respPatch.Wallet = wallet
	}

	if src.SellStrategies != nil {
		strats := mapDTOToStrategyConfigs(*src.SellStrategies)
		respPatch.SellStrategies = &strats
	}

	if src.TokenAddress != nil {
		token, err := solana.PublicKeyFromBase58(*src.TokenAddress)
		if err != nil {
			return nil, err
		}
		respPatch.Token = token.ToPointer()
	}

	return &respPatch, nil
}

func createSellPatch(src dto.TradingTaskPatch, wallet *wallets.SolanaWallet, rpcGroup *rpcgroups.RPCGroup) (strategies.Patch, error) {
	respPatch := strategies.SellPatch{}
	if src.SellAmount != nil {
		respPatch.SellAmount = src.SellAmount
	}

	if src.SellFee != nil {
		respPatch.SellFee = src.SellFee
	}

	if src.ComputeUnits != nil {
		respPatch.ComputeUnits = src.ComputeUnits
	}

	if src.Slippage != nil {
		respPatch.Slippage = src.Slippage
	}

	//if the wallet was set in the patch, we get it in the controller and pass it down
	//this is to stop mixing of services and keep them as top level as possible
	if wallet != nil {
		respPatch.Wallet = wallet
	}

	if rpcGroup != nil {
		respPatch.RPCGroup = rpcGroup
	}

	if src.TokenAddress != nil {
		token, err := solana.PublicKeyFromBase58(*src.TokenAddress)
		if err != nil {
			return nil, err
		}
		respPatch.Token = token.ToPointer()
	}

	return &respPatch, nil
}
