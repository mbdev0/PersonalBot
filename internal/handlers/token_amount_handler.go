package handlers

import (
	"math/big"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
)

func GetBuyTokenAmountFrom(buyInSolanaLamports big.Int, bondingCurve string) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	return bonding_curve_decoder.GetBuyTokenAmountFrom(buyInSolanaLamports, bondingCurve)
}
