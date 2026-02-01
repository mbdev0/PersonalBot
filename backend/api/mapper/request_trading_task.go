package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/monitoring/filters"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapTradingTaskDtoToTradingTask(src dto.TradingTask, wallet wallets.SolanaWallet) (strategies.Task, error) {
	switch src.Type {
	case dto.AFK:
		return mapAfkDtoToAfk(src, wallet)
	case dto.BUY:
		return mapBuyDtoToBuy(src, wallet)
	case dto.SELL:
		return mapSellDtoToSell(src, wallet)
	default:
		return nil, fmt.Errorf("task with type: %s - not found", src.Type)
	}
}

func mapBuyDtoToBuy(src dto.TradingTask, wallet wallets.SolanaWallet) (dst *strategies.Buy, err error) {
	dest := strategies.Buy{}
	dest.New()

	dest.BuyFee = *src.BuyFee
	dest.ComputeUnits = float64(src.ComputeUnits)
	dest.Slippage = src.Slippage
	dest.BuyAmount = utils.ConvertSolToLamport(*src.BuyAmount)
	dest.Wallet = wallet
	dest.SellStrategies = mapDTOToStrategyConfigs(*src.SellStrategies)
	dest.SellFee = *src.SellFee
	dest.State = string(strategies.CREATED)

	dest.Token, err = solana.PublicKeyFromBase58(*src.TokenAddress)
	if err != nil {
		return nil, err
	}

	return &dest, nil

}

func mapAfkDtoToAfk(src dto.TradingTask, wallet wallets.SolanaWallet) (dst *strategies.Afk, err error) {
	dest := strategies.Afk{}
	dest.New()

	dest.BuyFee = *src.BuyFee
	dest.ComputeUnits = float64(src.ComputeUnits)
	dest.Slippage = src.Slippage
	dest.BuyAmount = utils.ConvertSolToLamport(*src.BuyAmount)
	dest.Wallet = wallet
	dest.Filters = mapFiltersToDestFilters(*src.Filters)
	dest.SellStrategies = mapDTOToStrategyConfigs(*src.SellStrategies)
	dest.SellFee = *src.SellFee
	dest.State = string(strategies.CREATED)

	return &dest, nil
}

func mapSellDtoToSell(src dto.TradingTask, wallet wallets.SolanaWallet) (dst strategies.Task, err error) {
	dest := strategies.Sell{}
	dest.New()

	dest.SellFee = *src.SellFee
	dest.SellAmount = *src.SellAmount
	dest.ComputeUnits = float64(src.ComputeUnits)
	dest.Slippage = src.Slippage
	dest.Wallet = wallet
	dest.State = string(strategies.CREATED)

	dest.Token, err = solana.PublicKeyFromBase58(*src.TokenAddress)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func mapDTOToStrategyConfigs(dtos []dto.SellStrategyDTO) []strategies.StrategyConfig {
	configs := make([]strategies.StrategyConfig, len(dtos))
	for i, dto := range dtos {
		configs[i] = strategies.StrategyConfig{
			Type:       strategies.SellStrategyType(dto.Type),
			Value:      dto.Value,
			SellAmount: dto.SellAmount,
		}
	}
	return configs
}

func mapFiltersToDestFilters(srcFilters dto.Filters) []strategies.StrategyFilter {
	destFilters := make([]strategies.StrategyFilter, 0)

	extractAndCheck := func(f *bool) bool {
		return f != nil && *f
	}

	if extractAndCheck(srcFilters.HasTelegram) {
		destFilters = append(destFilters, filters.HasTelegram)
	}
	if extractAndCheck(srcFilters.HasTwitter) {
		destFilters = append(destFilters, filters.HasTwitter)
	}
	if extractAndCheck(srcFilters.HasWebsite) {
		destFilters = append(destFilters, filters.HasWebsite)
	}

	if srcFilters.DevWallet != nil && *srcFilters.DevWallet != "" {
		wallet := *srcFilters.DevWallet
		destFilters = append(destFilters, func() filters.FilterInfo {
			return filters.DevWallet(wallet)
		})
	}

	return destFilters
}
