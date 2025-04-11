package handlers

import (
	"math/big"
	buy_utils "pump_fun/internal/buy/utils.go"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
)

func GetBuyTokenAmountFrom(buyInSolanaLamports big.Int, bondingCurve string) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	return bonding_curve_decoder.GetBuyTokenAmountFrom(buyInSolanaLamports, bondingCurve)
}

func AddSlippageToBuy(lamportAmount big.Int, slippagePercentage float64) (newBuyAmount big.Int) {
	return buy_utils.AddSlippageToBuy(lamportAmount, slippagePercentage)
}

func ConvertSolToLamport(solAmount float64) (lamportAmount big.Int) {
	return buy_utils.ConvertSolToLamport(solAmount)
}
