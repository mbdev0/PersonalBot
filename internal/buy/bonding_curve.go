package buy

import (
	"errors"
	"fmt"
	"math/big"
	"pump_fun/internal/launch"
	"pump_fun/internal/models"
	rpccalls "pump_fun/internal/rpc_calls"

	"github.com/mr-tron/base58"
)

func GetMarketCapFrom(bondingCurveValue string) (marketCapVal *big.Float, err error, hasCompleted bool) {
	bondingCurveData, err, hasCompleted := getBondingCurveData(bondingCurveValue)
	if err != nil {
		return nil, err, false
	}

	if hasCompleted {
		return nil, fmt.Errorf("coin has already migrated"), true
	}

	marketCap, err := getMarketCap(*bondingCurveData)
	if err != nil {
		return nil, err, false
	}

	return marketCap, nil, false
}

func GetMarketCapInitial(bondingCurveAddress string) (marketCapVal *big.Float, err error, hasCompleted bool) {
	bondingCurveResponse, err := rpccalls.GetAccountInfo(bondingCurveAddress)
	if err != nil {
		return nil, err, false
	}

	decodedBondingCurve := bondingCurveResponse.Value.Data.GetBinary()
	value := base58.Encode(decodedBondingCurve)

	bondingCurveData, err, hasCompleted := getBondingCurveData(value)
	if err != nil {
		return nil, err, false
	}

	if hasCompleted {
		return nil, fmt.Errorf("coin has already migrated"), true
	}

	marketCap, err := getMarketCap(*bondingCurveData)
	if err != nil {
		return nil, err, false
	}

	return marketCap, nil, false

}

func GetBuyTokenAmountFrom(buyInSol big.Int, bondingCurve string) (tokenAmnt *big.Int, err error, hasCompleted bool) {

	bondingCurveResponse, err := rpccalls.GetAccountInfo(bondingCurve)
	if err != nil {
		return nil, err, false
	}

	decodedBondingCurve := bondingCurveResponse.Value.Data.GetBinary()
	value := base58.Encode(decodedBondingCurve)

	bondingCurveData, err, hasCompleted := getBondingCurveData(value)
	if err != nil {
		return nil, err, false
	}

	if hasCompleted {
		return nil, fmt.Errorf("coin has already migrated"), true
	}

	tokenAmount := getTokenAmount(buyInSol, bondingCurveData)
	return &tokenAmount, nil, false
}

func getBondingCurveData(bondingCurve string) (bondingCurveData *models.BondingCurve, err error, hasMigrated bool) {

	bondingCurveDataBytes, err := base58.Decode(bondingCurve)
	if err != nil {
		return nil, err, false
	}

	if bondingCurveDataBytes[len(bondingCurveDataBytes)-1] == 1 {
		return nil, nil, true
	}

	bondingCurveInfo, err := decryptBondingCurveData(bondingCurveDataBytes)
	if err != nil {
		return nil, err, false
	}

	return bondingCurveInfo, nil, false
}

func getMarketCap(bondingCurve models.BondingCurve) (*big.Float, error) {
	// This has been verified as accurate comparing against pump.funs website
	solPrice, err := launch.GetSolPrice() // -> the 100 will be gotten from an api call that runs periodically in the background https://frontend-api-v3.pump.fun/sol-price
	if err != nil {
		return nil, err
	}

	bigSolPrice := new(big.Float).SetFloat64(*solPrice)
	floatSolRes := new(big.Float).SetInt(&bondingCurve.VirtualSolReserves)
	floatTokenReserves := new(big.Float).SetInt(&bondingCurve.VirtualTokenReserves)
	marketCapInSol := new(big.Float).Quo(floatSolRes, floatTokenReserves)
	marketCap := new(big.Float).Mul(marketCapInSol, bigSolPrice)

	return new(big.Float).Mul(marketCap, big.NewFloat(1000000)), nil
}

func getTokenAmount(solBuy big.Int, bondingCurveInfo *models.BondingCurve) big.Int {
	i := new(big.Int).Add(&bondingCurveInfo.VirtualSolReserves, &solBuy)
	n := new(big.Int).Mul(&bondingCurveInfo.VirtualSolReserves, &bondingCurveInfo.VirtualTokenReserves)
	r := new(big.Int).Div(n, i)
	s := new(big.Int).Sub(&bondingCurveInfo.VirtualTokenReserves, r)

	if s.Cmp(&bondingCurveInfo.RealTokenReserves) == -1 {
		return *s
	}

	return bondingCurveInfo.RealTokenReserves
}

func decryptBondingCurveData(dataBinary []byte) (*models.BondingCurve, error) {
	if len(dataBinary) != 49 {
		return nil, errors.New("base58 string for bonding curve is too short")
	}

	bondingCurve := models.BondingCurve{}
	bondingCurve.VirtualTokenReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[8:16]))
	bondingCurve.VirtualSolReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[16:24]))
	bondingCurve.RealTokenReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[24:32]))
	bondingCurve.RealSolReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[32:40]))
	bondingCurve.MaxTokens = *new(big.Int).SetBytes(reverseBytes(dataBinary[40:48]))
	bondingCurve.IsCompleted = dataBinary[48] != 0

	return &bondingCurve, nil
}

func reverseBytes(b []byte) []byte {
	reversed := make([]byte, len(b))
	for i := range b {
		reversed[len(b)-1-i] = b[i]
	}
	return reversed
}
