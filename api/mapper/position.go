package mapper

import (
	"pump_fun/api/dto"
	"pump_fun/internal/core/position"
)

func MapPositionToPositionDto(src position.Position) (dest dto.PositionDto) {
	dest.PositionId = src.PositionId
	dest.FinalizedProfit = src.FinalizedProfit.Text('f', 9)
	dest.InitialTokenAmount = src.InitialTokenAmount.Text('f', 9)
	dest.RemaningCostBasis = src.InitialTokenAmount.Text('f', 9)
	dest.TokenAddress = src.TokenAddress.String()
	dest.WalletAddress = src.WalletAddress.String()
	dest.TokenRemaining = src.TokenRemaining.Text('f', 9)

	return dest
}
