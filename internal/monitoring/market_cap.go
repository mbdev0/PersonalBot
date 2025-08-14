package monitoring

import (
	"context"
	"pump_fun/internal/monitoring/stream"
	"pump_fun/internal/monitoring/stream/response"
	bondingcurve "pump_fun/internal/solana/programs/pumpfun/bonding_curve"
	"pump_fun/pkg/logger"
)

/*

Example usage:

	ctx, cancel := context.WithCancel(context.Background()) //we create a context

	bondingCurveAddress := "5U3smn2USzQGSWfM3JmKXt9YPKpWifjxvR2aMZkQAN1S"
	go monitoring.StartMarketCapMonitor(ctx, bondingCurveAddress) //starts monitoring of marketcap of a certain bonding address (we give bonding adress rather than CA)

	tokenAmount, err, hasCompleted := buy.GetBuyTokenAmountFrom(*big.NewInt(10000000), bondingCurveAddress) //get's token amount for a certain bondingCurveAddress, we use lamports for solana
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenAmount.String())

	time.Sleep(time.Second * 10)
	cancel() // this will cancel the context and thus close down the go routine from outside.. sick yeah?

*/

func StartMarketCapMonitor(ctx context.Context, bondingCurveAddress string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		marketCapMonitor(ctx, bondingCurveAddress)
	}
}

func marketCapMonitor(ctx context.Context, bondingCurveAddress string) {
	marketCapInit, err, hasCompleted := bondingcurve.GetMarketCapInitial(bondingCurveAddress, ctx)

	if err != nil {
		logger.Error("Error getting initial market cap", err)
		if hasCompleted {
			return
		}
	}

	logger.Information("Initial market cap: ", marketCapInit.String())

	accountInfoChan := make(chan response.AccountSubscribeModel, 20)

	go stream.Geyser_Stream_AccountInfo(ctx, bondingCurveAddress, accountInfoChan)

	for accountInfo := range accountInfoChan {
		marketCap, err, hasCompleted := bondingcurve.GetMarketCapFrom(accountInfo.Params.Result.Value.Data[0])
		if err != nil {
			logger.Error("Error getting market cap", err)
			if hasCompleted {
				return
			}
		}
		logger.Information("Market cap: ", marketCap.String())
	}

}
