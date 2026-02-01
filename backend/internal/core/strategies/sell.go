package strategies

import (
	"personal_bot/internal/core/models/wallets"

	"github.com/gagliardetto/solana-go"
)

type Sell struct {
	id           int64
	strategyType TradingType
	SellAmount   float64
	SellFee      float64
	Slippage     float64
	ComputeUnits float64
	Token        solana.PublicKey
	Wallet       wallets.SolanaWallet
	State        string
	SellTaskId   int64
}

func (b *Sell) New() {
	b.strategyType = SELL
}

func (b *Sell) StrategyTaskId() int64 {
	return b.id
}

func (b *Sell) SetId(id int64) {
	b.id = id
}
func (b *Sell) StrategyType() TradingType {
	return b.strategyType
}

func (b *Sell) StrategyState() string {
	return b.State
}

func (b *Sell) SetStrategyState(st string) {
	b.State = st
}

func (b *Sell) GetWallet() wallets.SolanaWallet {
	return b.Wallet
}
func (b *Sell) GetComputeUnits() uint32 {
	return uint32(b.ComputeUnits)
}
func (b *Sell) GetSlippage() float64 {
	return b.Slippage
}
func (b *Sell) GetSellFee() float64 {
	return b.SellFee
}
