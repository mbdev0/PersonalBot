package monitoring

import (
	"context"
	"math/big"
)

type MarketCapMonitor interface {
	StreamMarketCap(ctx context.Context) chan *big.Float
}
