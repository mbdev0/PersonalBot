package utils

import (
	"math/big"
	"pump_fun/internal/constants"
)

func ConvertSolToLamport(solAmount float64) (lamportAmount big.Int) {
	solAmountBigFloat := new(big.Float).SetFloat64(solAmount)
	lamportConversion := new(big.Float).SetInt64(constants.LamportsConversion)
	floatLamports := new(big.Float).Mul(solAmountBigFloat, lamportConversion)
	intLamports, _ := floatLamports.Int(new(big.Int))

	return *intLamports

}
