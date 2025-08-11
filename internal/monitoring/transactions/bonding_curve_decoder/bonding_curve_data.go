package bonding_curve_decoder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"pump_fun/internal/core/models"
	"pump_fun/internal/launch/solana_price"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
)

var marketCap_Multiplier float64 = 1000000

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

func GetMarketCapInitial(bondingCurveAddress string, cancellationToken models.CancelToken) (marketCapVal *big.Float, err error, hasCompleted bool) {
	bondingCurveResponse, err := rpcclient.GetAccountInfo(bondingCurveAddress, cancellationToken)
	if err != nil {
		return nil, err, false
	}

	decodedBondingCurve := bondingCurveResponse.Value.Data.GetBinary()
	value := base64.StdEncoding.EncodeToString(decodedBondingCurve)

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

func getMarketCap(bondingCurve models.BondingCurve) (*big.Float, error) {
	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		return nil, err
	}

	bigSolPrice := new(big.Float).SetFloat64(*solPrice)
	floatSolRes := new(big.Float).SetInt(&bondingCurve.VirtualSolReserves)
	floatTokenReserves := new(big.Float).SetInt(&bondingCurve.VirtualTokenReserves)
	marketCapInSol := new(big.Float).Quo(floatSolRes, floatTokenReserves)
	marketCap := new(big.Float).Mul(marketCapInSol, bigSolPrice)

	return new(big.Float).Mul(marketCap, big.NewFloat(marketCap_Multiplier)), nil
}

func GetSolanaTokenPrice(bondingCurve models.BondingCurve, tokenAmount uint64) *uint64 {
	floatSolRes := new(big.Float).SetInt(&bondingCurve.VirtualSolReserves)
	floatTokenReserves := new(big.Float).SetInt(&bondingCurve.VirtualTokenReserves)
	marketCapInSol := new(big.Float).Quo(floatSolRes, floatTokenReserves)
	tokenPrice, _ := new(big.Float).Mul(marketCapInSol, big.NewFloat(float64(tokenAmount))).Uint64()

	return &tokenPrice

}

func GetBondingCurveDataFromAddress(bondingCurveAddress string, cancellationToken models.CancelToken) (bondingCurveData *models.BondingCurve, err error, hasCompleted bool) {
	bondingCurveResponse, err := rpcclient.GetAccountInfo(bondingCurveAddress, cancellationToken)
	if err != nil {
		return nil, err, false
	}

	decodedBondingCurve := bondingCurveResponse.Value.Data.GetBinary()
	value := base64.StdEncoding.EncodeToString(decodedBondingCurve)

	bondingCurveModel, err, hasCompleted := getBondingCurveData(value)
	if err != nil {
		return nil, err, false
	}

	if hasCompleted {
		return nil, fmt.Errorf("coin has already migrated"), true
	}

	return bondingCurveModel, nil, false
}

func GetBuyTokenAmountFrom(buyInSol big.Int, bondingCurveData *models.BondingCurve) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	tokenAmount := getTokenAmount(buyInSol, bondingCurveData)
	return &tokenAmount, nil, false
}

func getBondingCurveData(bondingCurve string) (bondingCurveData *models.BondingCurve, err error, hasMigrated bool) {

	bondingCurveDataBytes, err := base64.RawStdEncoding.DecodeString(bondingCurve)
	if err != nil {
		return nil, err, false
	}

	bondingCurveInfo, err := decryptBondingCurveData(bondingCurveDataBytes)
	if err != nil {
		return nil, err, false
	}

	if bondingCurveInfo.IsCompleted {
		return nil, nil, true
	}

	return bondingCurveInfo, nil, false
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
	if len(dataBinary) > 150 {
		return nil, errors.New("base64 string for bonding curve is too long")
	}

	bondingCurve := models.BondingCurve{}
	bondingCurve.VirtualTokenReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[8:16]))
	bondingCurve.VirtualSolReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[16:24]))
	bondingCurve.RealTokenReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[24:32]))
	bondingCurve.RealSolReserves = *new(big.Int).SetBytes(reverseBytes(dataBinary[32:40]))
	bondingCurve.MaxTokens = *new(big.Int).SetBytes(reverseBytes(dataBinary[40:48]))
	bondingCurve.IsCompleted = dataBinary[48] != 0
	bondingCurve.DevWallet = solana.PublicKeyFromBytes(dataBinary[49:81])

	return &bondingCurve, nil
}

func reverseBytes(b []byte) []byte {
	reversed := make([]byte, len(b))
	for i := range b {
		reversed[len(b)-1-i] = b[i]
	}
	return reversed
}
