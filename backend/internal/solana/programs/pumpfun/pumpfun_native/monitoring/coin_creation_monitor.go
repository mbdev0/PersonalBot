package monitoring

import (
	"context"
	"fmt"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/solana/monitoring/decoder"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/monitoring/models"
	"personal_bot/internal/solana/monitoring/stream"
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
}

func NewPumpfunNativeCoinCreation(wsUrl string, httpClient *rpc.Client, filters filters.FilterPipeline) *PumpfunNativeCoinCreation {
	return &PumpfunNativeCoinCreation{
		monitoringAddress: constants.PumpFunProgram,
		wsUrl:             wsUrl,
		httpClient:        httpClient,
		filters:           filters,
	}
}

func (pn *PumpfunNativeCoinCreation) StreamCoinCreation(ctx context.Context) chan models.Coin {
	return streamUtils.Stream[models.Coin](ctx, pn.monitoringAddress, pn.wsUrl, stream.NewStartGeyserTransactionStream, pn.parseTransactionForCoin)
}

func (pn *PumpfunNativeCoinCreation) parseTransactionForCoin(data []byte) (models.Coin, error) {
	tx, err := jsonparse.Decode[response.TransactionNotification](data)
	if err != nil {
		return models.Coin{}, err
	}

	coin := decoder.DecryptTransactionNotificationForCoin(tx, pn.filters.ShouldAccessIPFS())
	if coin != nil {
		coin = pn.filters.ApplyFilters(coin)
	}

	return models.Coin{}, fmt.Errorf("coin did not pass filters")
}
