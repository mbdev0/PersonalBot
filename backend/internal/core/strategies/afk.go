package strategies

import (
	"math/big"
	"pump_fun/app/iterable"
	"pump_fun/internal/core/models/wallets"
	"pump_fun/internal/monitoring/filters"
)

type StrategyFilter func() filters.FilterInfo

type Afk struct {
	id             int64
	strategyType   TradingType
	Filters        []StrategyFilter
	BuyAmount      *big.Int
	BuyFee         float64
	Slippage       float64
	ComputeUnits   float64
	Wallet         wallets.SolanaWallet
	SellStrategies []StrategyConfig
	SellFee        float64
}

func (a *Afk) New() {
	a.id = iterable.Itr.ID()
	a.strategyType = AFK
}

func (a *Afk) StrategyTaskId() int64 {
	return a.id
}

func (a *Afk) StrategyType() TradingType {
	return a.strategyType
}
