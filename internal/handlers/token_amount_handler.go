package handlers

import (
	"math/big"
	"pump_fun/internal/core/models"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
	"pump_fun/internal/utils"
)

func GetBuyTokenAmountFrom(buyInSolanaLamports big.Int, bondingCurveData *models.BondingCurve) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	return bonding_curve_decoder.GetBuyTokenAmountFrom(buyInSolanaLamports, bondingCurveData)
}

func ConvertSolToLamport(solAmount float64) (lamportAmount big.Int) {
	return utils.ConvertSolToLamport(solAmount)
}
