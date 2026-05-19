package monitoring

import (
	"context"
	"personal_bot/internal/solana/monitoring/models"
)

type CoinMonitor interface {
	StreamCoinCreations(ctx context.Context) chan models.Coin
}
