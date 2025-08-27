package strategies

import (
	"context"
	"math/big"
	"pump_fun/internal/monitoring/filters"

	"github.com/gagliardetto/solana-go"
)

type Afk struct {
	Filters      []func() filters.Filter
	BuyAmount    *big.Int
	BuyFee       float64
	Slippage     float64
	ComputeUnits float64
	Wallet       solana.PrivateKey
	Cancel       context.CancelFunc
}
