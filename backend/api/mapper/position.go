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

func MapPositionsToDashboard(src []position.Position) dto.PositionDashboard {
	//first we group
	grouped := map[string][]position.Position{}
	for _, pos := range src {
		key := pos.TokenAddress.String()
		grouped[key] = append(grouped[key], pos)
	}

	//then we loop through the values
	//make rows for each loop through the values
	resp := dto.PositionDashboard{}
	for coin, posGroup := range grouped {
		totalMarketCapEntry := big.NewFloat(0)
		totalMarketCapExit := big.NewFloat(0)
		totalFinalizedProfit := big.NewFloat(0)
		validMarketCapEntry := 0
		validMarketCapExit := 0

		for _, pos := range posGroup {
			//if greater than 0
			if pos.AverageMarketCapExit.Cmp(big.NewFloat(0)) == 1 {
				validMarketCapExit++
			}

			if pos.MarketCapEntry.Cmp(big.NewFloat(0)) == 1 {
				validMarketCapEntry++
			}
			totalMarketCapEntry.Add(totalMarketCapEntry, pos.MarketCapEntry)
			totalMarketCapExit.Add(totalMarketCapExit, pos.AverageMarketCapExit)
			totalFinalizedProfit.Add(totalFinalizedProfit, pos.FinalizedProfit)
		}

		if validMarketCapEntry == 0 {
			validMarketCapEntry = 1
		}
		if validMarketCapExit == 0 {
			validMarketCapExit = 1
		}

		averageMarketCapEntry := new(big.Float).Quo(totalMarketCapEntry, big.NewFloat(float64(validMarketCapEntry)))
		averageMarketCapExit := new(big.Float).Quo(totalMarketCapExit, big.NewFloat(float64(validMarketCapExit)))
		finalizedProfit := new(big.Float).Quo(totalFinalizedProfit, big.NewFloat(constants.LamportsConversion))

		posRow := dto.PositionDashboardRow{
			TotalPNL:              finalizedProfit.Text('f', 2),
			AverageMarketCapEntry: averageMarketCapEntry.Text('f', 2),
			AverageMarketCapExit:  averageMarketCapExit.Text('f', 2),
			Coin:                  coin,
		}

		if len(posGroup) > 0 {
			posRow.AddressForUrl = posGroup[0].AddressForUrl
		}

		resp = append(resp, posRow)

	}

	return resp
}
