package buy_utils

import (
	"math/big"
	"pump_fun/internal/constants"
)

func AddSlippageToBuy(lamportAmount big.Int, slippagePercentage float64) (newBuyAmount big.Int) {
	slippageFloat := new(big.Float).SetFloat64(slippagePercentage)
	slippageFloat = new(big.Float).Add(slippageFloat, big.NewFloat(1))

	lamportFloat := new(big.Float).SetInt(&lamportAmount)

	newLamportAmountFloat := new(big.Float).Mul(slippageFloat, lamportFloat)
	newLamportAmountInt, _ := newLamportAmountFloat.Int(new(big.Int))

	return *newLamportAmountInt
}

func ConvertSolToLamport(solAmount float64) (lamportAmount big.Int) {
	solAmountBigFloat := new(big.Float).SetFloat64(solAmount)
	lamportConversion := new(big.Float).SetInt64(constants.LamportsConversion)
	floatLamports := new(big.Float).Mul(solAmountBigFloat, lamportConversion)
	intLamports, _ := floatLamports.Int(new(big.Int))

	return *intLamports

}
