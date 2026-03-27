package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/solana/utils"

	"github.com/google/uuid"
)

func MapTradingTaskToDto(src strategies.Task) (dest *dto.TradingTaskResponse, err error) {
	switch t := src.(type) {
	case *strategies.Afk:
		return mapAfkToAfkResponse(t)
	case *strategies.Buy:
		return mapBuyToBuyResponse(t)
	case *strategies.Sell:
		return mapSellToSellResponse(t)
	default:
		return nil, fmt.Errorf("task not found with type - make sure the type is created/mapping is created")
	}
}

func mapBuyToBuyResponse(src *strategies.Buy) (dest *dto.TradingTaskResponse, err error) {
	dst := dto.TradingTaskResponse{}
	dst.Type = dto.TradingType(src.StrategyType())
	dst.Id = src.StrategyTaskId()
	buyAmount := utils.ConvertLamportToSol(src.BuyAmount)
	dst.BuyAmount = &buyAmount
	dst.BuyFee = &src.BuyFee
	dst.ComputeUnits = src.ComputeUnits
	dst.Slippage = src.Slippage
	dst.SellFee = src.SellFee
	sellStrategies := mapSellStratsToResponseStrats(src.SellStrategies)
	dst.SellStrategies = &sellStrategies
	dst.WalletName = src.Wallet.WalletName
	dst.WalletAddress = src.Wallet.PublicKey.Short(constants.ShortPublicAddressInt)
	dst.State = string(src.State)
	dst.BuyTaskId = &src.BuyTaskId
	dst.PositionId = &src.PositionId
	tokenAddress := src.Token.String()
	dst.TokenAddress = &tokenAddress
	dst.TimeCreated = src.TimeCreated
	dst.Message = src.Message
	dst.RPCGroup = src.RPCGroup.Name

	return &dst, nil
}

func mapAfkToAfkResponse(src *strategies.Afk) (dest *dto.TradingTaskResponse, err error) {
	dst := dto.TradingTaskResponse{}
	dst.Type = dto.TradingType(src.StrategyType())
	dst.Id = src.StrategyTaskId()
	buyAmount := utils.ConvertLamportToSol(src.BuyAmount)
	dst.BuyAmount = &buyAmount

	dst.BuyFee = &src.BuyFee
	dst.ComputeUnits = src.ComputeUnits
	dst.Slippage = src.Slippage

	filters := mapFiltersToResponseFilters(src.Filters)
	dst.Filters = &filters
	dst.SellFee = src.SellFee

	sellStrategies := mapSellStratsToResponseStrats(src.SellStrategies)
	dst.SellStrategies = &sellStrategies

	dst.WalletName = src.Wallet.WalletName
	dst.WalletAddress = src.Wallet.PublicKey.Short(constants.ShortPublicAddressInt)
	dst.State = string(src.State)
	dst.TimeCreated = src.TimeCreated
	dst.Message = src.Message
	dst.RPCGroup = src.RPCGroup.Name

	return &dst, nil
}

func mapSellToSellResponse(t *strategies.Sell) (dest *dto.TradingTaskResponse, err error) {
	dst := dto.TradingTaskResponse{}
	dst.Type = dto.TradingType(t.StrategyType())
	dst.Id = t.StrategyTaskId()
	sellAmount := t.SellAmount
	dst.SellAmount = &sellAmount
	dst.SellFee = &t.SellFee
	dst.ComputeUnits = t.ComputeUnits
	dst.Slippage = t.Slippage
	dst.WalletName = t.Wallet.WalletName
	dst.WalletAddress = t.Wallet.PublicKey.Short(constants.ShortPublicAddressInt)
	dst.State = string(t.State)
	dst.SellTaskId = &t.SellTaskId
	tokenAddress := t.Token.String()
	dst.TokenAddress = &tokenAddress
	dst.TimeCreated = t.TimeCreated
	dst.Message = t.Message
	dst.RPCGroup = t.RPCGroup.Name

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
			Id:         uuid.NewString(),
			Type:       string(config.Type),
			Value:      config.Value,
			SellAmount: config.SellAmount,
		}
	}

	return dest
}
