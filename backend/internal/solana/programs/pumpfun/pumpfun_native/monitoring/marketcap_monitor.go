package monitoring

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	datastream "personal_bot/internal/solana/monitoring/data_stream"
	"personal_bot/internal/solana/monitoring/models"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	jsonparse "personal_bot/pkg/json_parse"
	"personal_bot/pkg/logger"
	streamUtils "personal_bot/pkg/stream"

	"github.com/gagliardetto/solana-go/rpc"
)

type PumpfunNativeMarketcap struct {
	monitoringAddress string
	httpClient        *rpc.Client
	wsUrl             string
	datastream        datastream.DataStream
}

func NewPumpfunNativeMarketcap(bondingCurveAddress, wsUrl string, httpClient *rpc.Client, datastream datastream.DataStream) *PumpfunNativeMarketcap {
	return &PumpfunNativeMarketcap{
		monitoringAddress: bondingCurveAddress,
		wsUrl:             wsUrl,
		httpClient:        httpClient,
		datastream:        datastream,
	}
}

func (pm *PumpfunNativeMarketcap) GetInitialMarketCap(ctx context.Context) (marketCapVal *big.Float, err error, hasCompleted bool) {
	bondingCurve, err, hasComplete := bondingcurve.GetBondingCurveDataFromAddress(ctx, pm.monitoringAddress, pm.httpClient)
	if err != nil {
		if hasComplete {
			return nil, nil, true
		}
		return nil, nil, false
	}

	if bondingCurve == nil {
		return nil, fmt.Errorf("bonding curve was nil"), false
	}

	marketCapVal, err = bondingcurve.GetMarketCap(*bondingCurve)
	if err != nil {
		return nil, err, false
	}

	return marketCapVal, nil, false
}

func (pm *PumpfunNativeMarketcap) GetMarketCapFrom(data []byte) (*big.Float, error) {

	accountInfo, err := jsonparse.Decode[models.AccountSubscribeModel](data)
	if err != nil {
		return nil, err
	}

	var bondingCurveData []byte
	if len(accountInfo.Params.Result.Value.Data) > 0 {
		bondingCurveData, err = base64.StdEncoding.DecodeString(accountInfo.Params.Result.Value.Data[0])
		if err != nil {
			return nil, err
		}
	}

	bondingCurve, err, _ := bondingcurve.GetBondingCurveData(bondingCurveData)
	if err != nil {
		return nil, err
	}

	if bondingCurve == nil {
		return nil, fmt.Errorf("bonding curve was nil")
	}

	marketCapVal, err := bondingcurve.GetMarketCap(*bondingCurve)
	if err != nil {
		return nil, err
	}

	return marketCapVal, nil
}

func (pm *PumpfunNativeMarketcap) StreamMarketCap(ctx context.Context) <-chan big.Float {
	cleanup := func() {
		pm.datastream.Unsubscribe(pm.wsUrl, pm.monitoringAddress)
	}

	output := streamUtils.Stream[big.Float](ctx, pm.monitoringAddress, pm.wsUrl, pm.datastream.SubscribeToAccountStream, pm.GetMarketCapFrom, cleanup)

	mcap, err, _ := pm.GetInitialMarketCap(ctx)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if mcap != nil {
		output <- *mcap
	}

	return output
}
