package monitoring

import (
	"context"
	"personal_bot/internal/core/constants"
	datastream "personal_bot/internal/solana/monitoring/data_stream"
	"personal_bot/internal/solana/monitoring/decoder"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/monitoring/models"
	"personal_bot/internal/solana/monitoring/stream/response"
	jsonparse "personal_bot/pkg/json_parse"
	streamUtils "personal_bot/pkg/stream"

	"github.com/gagliardetto/solana-go/rpc"
)

type PumpfunNativeCoinCreation struct {
	monitoringAddress string
	httpClient        *rpc.Client
	wsUrl             string
	filters           filters.FilterPipeline
	datastream        datastream.DataStream
}

func NewPumpfunNativeCoinCreation(wsUrl string, httpClient *rpc.Client, filters filters.FilterPipeline, datastream datastream.DataStream) *PumpfunNativeCoinCreation {
	return &PumpfunNativeCoinCreation{
		monitoringAddress: constants.PumpFunProgram,
		wsUrl:             wsUrl,
		httpClient:        httpClient,
		filters:           filters,
		datastream:        datastream,
	}
}

func (pn *PumpfunNativeCoinCreation) StreamCoinCreations(ctx context.Context) chan models.Coin {
	cleanup := func() {
		pn.datastream.Unsubscribe(pn.wsUrl, pn.monitoringAddress)
	}

	return streamUtils.Stream[models.Coin](ctx, pn.monitoringAddress, pn.wsUrl, pn.datastream.SubscribeToTransaction, pn.parseTransactionForCoin, cleanup)
}

func (pn *PumpfunNativeCoinCreation) parseTransactionForCoin(data []byte) (*models.Coin, error) {
	tx, err := jsonparse.Decode[response.TransactionNotification](data)
	if err != nil {
		return nil, err
	}

	coin := decoder.DecryptTransactionNotificationForCoin(tx, pn.filters.ShouldAccessIPFS())
	if coin != nil {
		coin = pn.filters.ApplyFilters(coin)
	}
	return coin, nil
}
