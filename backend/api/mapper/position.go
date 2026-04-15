package mapper

import (
	"math/big"
	"personal_bot/api/dto"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/position"
)

func MapPositionToPositionDto(src position.Position) (dest dto.PositionDto) {
	dest.PositionId = src.PositionId

	finalizedProfit := new(big.Float).Quo(src.FinalizedProfit, big.NewFloat(constants.LamportsConversion))
	initialTokenAmount := new(big.Float).Quo(src.InitialTokenAmount, big.NewFloat(constants.TokenAmountDecimals))
	remainingCostBasis := new(big.Float).Quo(src.RemainingCostBasis, big.NewFloat(constants.LamportsConversion))
	tokensRemaining := new(big.Float).Quo(src.TokenRemaining, big.NewFloat(constants.TokenAmountDecimals))

	dest.FinalizedProfit = finalizedProfit.Text('f', 9)
	dest.InitialTokenAmount = initialTokenAmount.Text('f', 9)
	dest.RemaningCostBasis = remainingCostBasis.Text('f', 9)
	dest.TokenAddress = src.TokenAddress.String()
	dest.WalletAddress = src.WalletAddress.String()
	dest.TokenRemaining = tokensRemaining.Text('f', 9)
	dest.MarketCapEntry = src.MarketCapEntry.Text('f', 9)
	dest.AverageMarketCapExit = src.AverageMarketCapExit.Text('f', 9)

	return dest
}
