package monitoring

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/internal/solana/monitoring/stream"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	"personal_bot/pkg/logger"
	streamUtils "personal_bot/pkg/stream"

	"github.com/gagliardetto/solana-go/rpc"
)

type PumpfunNativeMarketcap struct {
	monitoringAddress string
	httpClient        *rpc.Client
	wsUrl             string
}

func NewPumpfunNativeMarketcap(bondingCurveAddress, wsUrl string, httpClient *rpc.Client) *PumpfunNativeMarketcap {
	return &PumpfunNativeMarketcap{
		monitoringAddress: bondingCurveAddress,
		wsUrl:             wsUrl,
		httpClient:        httpClient,
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

func (pm *PumpfunNativeMarketcap) GetMarketCapFrom(data []byte) (marketCapVal *big.Float, err error) {
	bondingCurve, err, _ := bondingcurve.GetBondingCurveData(data)
	if err != nil {
		return nil, err
	}

	if bondingCurve == nil {
		return nil, fmt.Errorf("bonding curve was nil")
	}

	marketCapVal, err = bondingcurve.GetMarketCap(*bondingCurve)
	if err != nil {
		return nil, err
	}

	return marketCapVal, nil
}

func (pm *PumpfunNativeMarketcap) StreamMarketCap(ctx context.Context) chan *big.Float {
	output := streamUtils.Stream[*big.Float](ctx, pm.monitoringAddress, pm.wsUrl, stream.NewGeyserStreamAccountInfo, pm.GetMarketCapFrom)

	mcap, err, _ := pm.GetInitialMarketCap(ctx)
	if err != nil {
		logger.Error(err)
		return nil
	}

	output <- mcap

	return output
}
