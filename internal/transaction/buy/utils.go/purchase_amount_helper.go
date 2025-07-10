package buy_utils

import (
	"math/big"
)

func AddSlippageToBuy(lamportAmount big.Int, slippagePercentage float64) (newBuyAmount big.Int) {
	slippageFloat := new(big.Float).SetFloat64(slippagePercentage)
	slippageFloat = new(big.Float).Add(slippageFloat, big.NewFloat(1))

	lamportFloat := new(big.Float).SetInt(&lamportAmount)

	newLamportAmountFloat := new(big.Float).Mul(slippageFloat, lamportFloat)
	newLamportAmountInt, _ := newLamportAmountFloat.Int(new(big.Int))

	return *newLamportAmountInt
}
