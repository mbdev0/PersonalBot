package monitoring

import (
	"context"
	"personal_bot/backend/internal/core/constants"
	datastream "personal_bot/backend/internal/solana/monitoring/data_stream"
	"personal_bot/backend/internal/solana/monitoring/filters"
	"personal_bot/backend/internal/solana/monitoring/models"
	transactiondecoder "personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/transaction_decoder"

	"personal_bot/backend/pkg/stream"

	"github.com/gagliardetto/solana-go/rpc"
)

type PumpfunAMMCoinCreationMonitor struct {
	datastream        datastream.DataStream
	monitoringAddress string
	wsUrl             string
	httpClient        *rpc.Client
	filters           filters.FilterPipeline
}

func NewPumpfunAMMCoinCreationMonitor(ds datastream.DataStream, wsUrl string, httpClient *rpc.Client, filters filters.FilterPipeline) *PumpfunAMMCoinCreationMonitor {
	return &PumpfunAMMCoinCreationMonitor{
		datastream:        ds,
		monitoringAddress: constants.PumpfunMigration,
		wsUrl:             wsUrl,
		httpClient:        httpClient,
		filters:           filters,
	}
}

func (paccm *PumpfunAMMCoinCreationMonitor) StreamCoinCreation(ctx context.Context) chan models.Coin {
	parser := func(b []byte) (*models.Coin, error) {
		return transactiondecoder.ParseTransactionNotification(ctx, b, paccm.httpClient, paccm.filters)
	}

	cleanup := func() {
		paccm.datastream.Unsubscribe(paccm.wsUrl, paccm.monitoringAddress)
	}

	return stream.Stream[models.Coin](ctx,
		paccm.monitoringAddress, paccm.wsUrl,
		paccm.datastream.SubscribeToTransaction,
		parser,
		cleanup,
	)
}
