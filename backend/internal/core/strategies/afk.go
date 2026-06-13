package strategies

import (
	"math/big"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/wallets"
	"personal_bot/internal/solana/monitoring/filters"
)

type StrategyFilter func() filters.FilterInfo

type Afk struct {
	Program        string
	id             int64
	strategyType   TradingType
	Filters        []StrategyFilter
	BuyAmount      *big.Int
	BuyFee         float64
	Slippage       float64
	ComputeUnits   float64
	Wallet         wallets.SolanaWallet
	SellStrategies []StrategyConfig
	SellFee        *float64
	State          string
	Message        string
	TimeCreated    int64
	RPCGroup       rpcgroups.RPCGroup
	Retries        *uint16
	RetriesDelayMS *uint32
}

func (a *Afk) New() {
	a.strategyType = AFK
}

func (a *Afk) GetProgram() string {
	return a.Program
}

func (a *Afk) StrategyTaskId() int64 {
	return a.id
}

func (a *Afk) SetId(id int64) {
	a.id = id
}

func (a *Afk) StrategyType() TradingType {
	return a.strategyType
}

func (a *Afk) StrategyState() string {
	return a.State
}

func (a *Afk) SetStrategyState(st string) {
	a.State = st
}

func (a *Afk) StrategyMessage() string {
	return a.Message
}

func (a *Afk) RPCGroupId() int64 {
	return a.RPCGroup.Id
}

func (a *Afk) SetStrategyMessage(message string) {
	a.Message = message
}

func (a *Afk) GetWallet() wallets.SolanaWallet {
	return a.Wallet
}
func (a *Afk) GetComputeUnits() uint32 {
	return uint32(a.ComputeUnits)
}
func (a *Afk) GetSlippage() float64 {
	return a.Slippage
}
func (a *Afk) GetSellFee() *float64 {
	return a.SellFee
}

func (a *Afk) GetTimeCreated() int64 {
	return a.TimeCreated
}

func (a *Afk) SetTimeCreated(time int64) {
	a.TimeCreated = time
}
