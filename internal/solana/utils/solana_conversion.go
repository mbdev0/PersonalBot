package utils

import (
	"math/big"
	"pump_fun/internal/core/constants"
)

func ConvertSolToLamport(solAmount float64) (lamportAmount *big.Int) {
	solAmountBigFloat := new(big.Float).SetFloat64(solAmount)
	lamportConversion := new(big.Float).SetInt64(constants.LamportsConversion)
	floatLamports := new(big.Float).Mul(solAmountBigFloat, lamportConversion)
	intLamports, _ := floatLamports.Int(new(big.Int))

	return intLamports
}

func ConvertLamportToSol(lamportAmount *big.Int) float64 {
	lamportFloat := new(big.Float).SetInt(lamportAmount)
	lamportConversion := new(big.Float).SetInt64(constants.LamportsConversion)
	solAmountBig := new(big.Float).Quo(lamportFloat, lamportConversion)
	solAmount, _ := solAmountBig.Float64()

	return solAmount
}
