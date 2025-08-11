package handlers

import (
	"math/big"
	"pump_fun/internal/core/models"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
	buy_utils "pump_fun/internal/transaction/buy/utils.go"
	"pump_fun/internal/utils"
)

func GetBuyTokenAmountFrom(buyInSolanaLamports big.Int, bondingCurveData *models.BondingCurve) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	return bonding_curve_decoder.GetBuyTokenAmountFrom(buyInSolanaLamports, bondingCurveData)
}

func AddSlippageToBuy(lamportAmount big.Int, slippagePercentage float64) (newBuyAmount big.Int) {
	return buy_utils.AddSlippageToBuy(lamportAmount, slippagePercentage)
}

func ConvertSolToLamport(solAmount float64) (lamportAmount big.Int) {
	return utils.ConvertSolToLamport(solAmount)
}
