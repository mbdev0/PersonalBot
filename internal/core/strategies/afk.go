package strategies

import (
	"context"
	"math/big"
	"pump_fun/internal/monitoring/filters"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type StrategyFilter func() filters.FilterInfo

type Afk struct {
	id           string
	strategyType TradingType
	Filters      []StrategyFilter
	BuyAmount    *big.Int
	BuyFee       float64
	Slippage     float64
	ComputeUnits float64
	Wallet       solana.PrivateKey
	Cancel       context.CancelFunc
}

func (a *Afk) New() {
	a.id = uuid.NewString()
}

func (a *Afk) StrategyTaskId() string {
	return a.id
}

func (a *Afk) StrategyType() TradingType {
	return a.strategyType
}
