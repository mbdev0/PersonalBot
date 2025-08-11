package handlers

import (
	"math/big"
	"pump_fun/internal/core/models"
	bondingcurve "pump_fun/internal/solana/programs/pumpfun/bonding_curve"
	"pump_fun/internal/utils"
)

func GetBuyTokenAmountFrom(buyInSolanaLamports big.Int, bondingCurveData *models.BondingCurve) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	return bondingcurve.GetBuyTokenAmountFrom(buyInSolanaLamports, bondingCurveData)
}

func ConvertSolToLamport(solAmount float64) (lamportAmount big.Int) {
	return utils.ConvertSolToLamport(solAmount)
}
