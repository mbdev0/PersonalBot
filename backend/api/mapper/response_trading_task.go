package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/solana/utils"
)

func MapTradingTaskToDto(src strategies.Task) (dest *dto.TradingTaskResponse, err error) {
	switch t := src.(type) {
	case *strategies.Afk:
		task, err := mapAfkToAfkResponse(t)
		if err != nil {
			return nil, err
		}
		return task, nil
	default:
		return nil, fmt.Errorf("task not found with type - make sure the type is created")
	}
}

func mapAfkToAfkResponse(src *strategies.Afk) (dest *dto.TradingTaskResponse, err error) {
	dst := dto.TradingTaskResponse{}
	dst.Type = dto.TradingType(src.StrategyType())
	dst.Id = src.StrategyTaskId()
	dst.BuyAmount = utils.ConvertLamportToSol(src.BuyAmount)
	dst.BuyFee = src.BuyFee
	dst.ComputeUnits = src.ComputeUnits
	dst.Slippage = src.Slippage
	dst.Filters = mapFiltersToResponseFilters(src.Filters)
	dst.SellFee = src.SellFee
	dst.SellStrategies = mapSellStratsToResponseStrats(src.SellStrategies)
	dst.WalletName = src.Wallet.WalletName
	dst.WalletAddress = src.Wallet.PublicKey.Short(constants.ShortPublicAddressInt)

	return &dst, nil
}

func mapFiltersToResponseFilters(src []strategies.StrategyFilter) dto.Filters {
	dest := dto.Filters{}

	for _, filterFunc := range src {
		filter := filterFunc()
		b := true

		switch filter.Name {
		case "HasWebsite":
			dest.HasWebsite = &b
		case "HasTwitter":
			dest.HasTwitter = &b
		case "HasTelegram":
			dest.HasTelegram = &b
		case "DevWallet":
			if filter.Value != "" {
				dest.DevWallet = &filter.Value
			}
		}
	}

	return dest
}

func mapSellStratsToResponseStrats(src []strategies.StrategyConfig) []dto.SellStrategyDTO {
	dest := make([]dto.SellStrategyDTO, len(src))

	for i, config := range src {
		dest[i] = dto.SellStrategyDTO{
			Type:       string(config.Type),
			Value:      config.Value,
			SellAmount: config.SellAmount,
		}
	}

	return dest
}
