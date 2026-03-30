package bondingcurve

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/models"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func GetMarketCapFrom(bondingCurveValue string) (marketCapVal *big.Float, err error, hasCompleted bool) {
	//check if base64 here, if so convert to bytes then pass in
	bondingCurveBytes, err := base64.StdEncoding.DecodeString(bondingCurveValue)
	if err != nil {
		return nil, err, false
	}
	bondingCurveData, err, hasCompleted := getBondingCurveData(bondingCurveBytes)
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

func GetMarketCapFromTokenAddress(tokenAddrress solana.PublicKey, ctx context.Context, httpNode *rpc.Client) (marketCapVal *big.Float, err error, hasCompleted bool) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(tokenAddrress.String())
	if err != nil {
		return nil, err, false
	}

	return GetMarketCapInitial(bondingCurveAddress, ctx, httpNode)
}

func GetMarketCapInitial(bondingCurveAddress string, ctx context.Context, httpNode *rpc.Client) (marketCapVal *big.Float, err error, hasCompleted bool) {
	bondingCurveResponse, err := client.GetAccountInfo(bondingCurveAddress, ctx, httpNode)
	if err != nil {
		return nil, err, false
	}

	decodedBondingCurve := bondingCurveResponse.Value.Data.GetBinary()

	bondingCurveData, err, hasCompleted := getBondingCurveData(decodedBondingCurve)
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

	virtualSol := new(big.Float).SetUint64(bondingCurve.VirtualSolReserves)
	virtualSol.Quo(virtualSol, big.NewFloat(constants.LamportsConversion))
	virtualTokens := new(big.Float).SetUint64(bondingCurve.VirtualTokenReserves)
	virtualTokens.Quo(virtualTokens, big.NewFloat(constants.TokenAmountDecimals))

	pricePerTokenSol := new(big.Float).Quo(virtualSol, virtualTokens)

	totalSupply := big.NewFloat(1_000_000_000)

	marketCapSol := new(big.Float).Mul(pricePerTokenSol, totalSupply)

	solPriceFloat := big.NewFloat(*solPrice)
	marketCapUSD := new(big.Float).Mul(marketCapSol, solPriceFloat)

	return marketCapUSD, nil
}

func GetSolanaTokenPrice(bondingCurve models.BondingCurve, tokenAmount uint64) *uint64 {
	floatSolRes := new(big.Float).SetUint64(bondingCurve.VirtualSolReserves)
	floatTokenReserves := new(big.Float).SetUint64(bondingCurve.VirtualTokenReserves)
	marketCapInSol := new(big.Float).Quo(floatSolRes, floatTokenReserves)
	tokenPrice, _ := new(big.Float).Mul(marketCapInSol, big.NewFloat(float64(tokenAmount))).Uint64()

	return &tokenPrice

}

func GetBondingCurveDataFromAddress(bondingCurveAddress string, ctx context.Context, httpNode *rpc.Client) (bondingCurveData *models.BondingCurve, err error, hasCompleted bool) {
	bondingCurveResponse, err := client.GetAccountInfo(bondingCurveAddress, ctx, httpNode)
	if err != nil {
		return nil, err, false
	}

	decodedBondingCurve := bondingCurveResponse.Value.Data.GetBinary()

	bondingCurveModel, err, hasCompleted := getBondingCurveData(decodedBondingCurve)
	if err != nil {
		return nil, err, false
	}

	if bondingCurveModel.Complete {
		return nil, fmt.Errorf("coin has already migrated"), true
	}

	return bondingCurveModel, nil, false
}

func getBondingCurveData(bondingCurveBytes []byte) (bondingCurveData *models.BondingCurve, err error, hasMigrated bool) {

	bondingCurveInfo, err := decryptBondingCurveData(bondingCurveBytes)
	if err != nil {
		return nil, err, false
	}

	if bondingCurveInfo.Complete {
		return nil, nil, true
	}

	return bondingCurveInfo, nil, false
}

func GetBuyTokenAmountFrom(buyInSol big.Int, bondingCurveData *models.BondingCurve) (tokenAmnt *big.Int, err error, hasCompleted bool) {
	tokenAmount := getTokenAmount(buyInSol, bondingCurveData)
	return &tokenAmount, nil, false
}

func getTokenAmount(solBuy big.Int, bondingCurveInfo *models.BondingCurve) big.Int {
	virtualTokenReserves := new(big.Int).SetUint64(bondingCurveInfo.VirtualTokenReserves)
	virtualSolReserves := new(big.Int).SetUint64(bondingCurveInfo.VirtualSolReserves)
	realTokenReserves := new(big.Int).SetUint64(bondingCurveInfo.RealTokenReserves)

	i := new(big.Int).Add(virtualSolReserves, &solBuy)
	n := new(big.Int).Mul(virtualSolReserves, virtualTokenReserves)
	r := new(big.Int).Div(n, i)
	s := new(big.Int).Sub(virtualTokenReserves, r)

	if s.Cmp(realTokenReserves) == -1 {
		return *s
	}

	return *realTokenReserves
}

func decryptBondingCurveData(dataBinary []byte) (*models.BondingCurve, error) {

	if len(dataBinary) < 8 {
		return nil, nil
	}

	if !bytes.Equal(dataBinary[:8], constants.BondingCurveDiscriminator[:]) {
		return nil, fmt.Errorf("incorrect discriminator")
	}

	bondingCurve := new(models.BondingCurve)
	var err error
	switch len(dataBinary) {
	case 83:
		bondingCurve, err = parseV1(dataBinary)
		if err != nil {
			return nil, err
		}
	case 151:
		bondingCurve, err = parseV2(dataBinary)
		if err != nil {
			return nil, err
		}

	}

	return bondingCurve, nil
}

func parseV0(data []byte) (*models.BondingCurve, error) {
	v0 := new(models.BondingCurveV0)
	err := borsh.Deserialize(v0, data[8:])
	if err != nil {
		logger.Error("err whilst doing new way: ", err)
		return nil, err
	}

	return &models.BondingCurve{
		VirtualTokenReserves: v0.VirtualTokenReserves,
		VirtualSolReserves:   v0.VirtualSolReserves,
		RealTokenReserves:    v0.RealTokenReserves,
		RealSolReserves:      v0.RealSolReserves,
		TokenTotalSupply:     v0.TokenTotalSupply,
		Complete:             v0.Complete,
	}, nil
}

func parseV1(data []byte) (*models.BondingCurve, error) {
	v1 := new(models.BondingCurveV1)
	err := borsh.Deserialize(v1, data[8:])
	if err != nil {
		logger.Error("err whilst doing new way: ", err)
		return nil, err
	}

	return &models.BondingCurve{
		VirtualTokenReserves: v1.VirtualTokenReserves,
		VirtualSolReserves:   v1.VirtualSolReserves,
		RealTokenReserves:    v1.RealTokenReserves,
		RealSolReserves:      v1.RealSolReserves,
		TokenTotalSupply:     v1.TokenTotalSupply,
		Complete:             v1.Complete,
		Creator:              v1.Creator.ToPointer(),
	}, nil
}

func parseV2(data []byte) (*models.BondingCurve, error) {
	v2 := new(models.BondingCurveV2)
	err := borsh.Deserialize(v2, data[8:])
	if err != nil {
		logger.Error("err whilst doing new way: ", err)
		return nil, err
	}

	return &models.BondingCurve{
		VirtualTokenReserves: v2.VirtualTokenReserves,
		VirtualSolReserves:   v2.VirtualSolReserves,
		RealTokenReserves:    v2.RealTokenReserves,
		RealSolReserves:      v2.RealSolReserves,
		TokenTotalSupply:     v2.TokenTotalSupply,
		Complete:             v2.Complete,
		Creator:              v2.Creator.ToPointer(),
		IsMayhemMode:         v2.IsMayhemMode,
	}, nil
}
