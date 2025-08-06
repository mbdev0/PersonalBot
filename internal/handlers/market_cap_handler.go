package handlers

import (
	"math/big"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
)

func GetMarketCapFrom(bondingCurveValue string) (marketCapVal *big.Float, err error, hasCompleted bool) {
	return bonding_curve_decoder.GetMarketCapFrom(bondingCurveValue)
}

func GetMarketCapInitial(bondingCurveAddress string, cancellationToken models.CancelToken) (marketCapVal *big.Float, err error, hasCompleted bool) {
	return bonding_curve_decoder.GetMarketCapInitial(bondingCurveAddress, cancellationToken)
}
