package strategies

import (
	"math/big"
	"personal_bot/app/iterable"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/monitoring/filters"
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
