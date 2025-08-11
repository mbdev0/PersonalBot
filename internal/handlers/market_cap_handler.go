package handlers

import (
	"math/big"
	"pump_fun/internal/core/models"
	bondingcurve "pump_fun/internal/solana/programs/pumpfun/bonding_curve"
)

func GetMarketCapFrom(bondingCurveValue string) (marketCapVal *big.Float, err error, hasCompleted bool) {
	return bondingcurve.GetMarketCapFrom(bondingCurveValue)
}

func GetMarketCapInitial(bondingCurveAddress string, cancellationToken models.CancelToken) (marketCapVal *big.Float, err error, hasCompleted bool) {
	return bondingcurve.GetMarketCapInitial(bondingCurveAddress, cancellationToken)
}
