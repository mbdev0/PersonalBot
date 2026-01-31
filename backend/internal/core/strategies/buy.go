package strategies

import (
	"math/big"
	"personal_bot/internal/core/models/wallets"

	"github.com/gagliardetto/solana-go"
)

type Buy struct {
	id             int64
	strategyType   TradingType
	BuyAmount      *big.Int
	BuyFee         float64
	Slippage       float64
	ComputeUnits   float64
	Token          solana.PublicKey
	Wallet         wallets.SolanaWallet
	SellStrategies []StrategyConfig
	SellFee        float64
	BuyTaskId      int64
	PositionId     int64
	State          string
}

func (b *Buy) New() {
	b.strategyType = BUY
}

func (b *Buy) StrategyTaskId() int64 {
	return b.id
}

func (b *Buy) SetId(id int64) {
	b.id = id
}
func (b *Buy) StrategyType() TradingType {
	return b.strategyType
}

func (b *Buy) StrategyState() string {
	return b.State
}

func (b *Buy) SetStrategyState(st string) {
	b.State = st
}

func (b *Buy) GetWallet() wallets.SolanaWallet {
	return b.Wallet
}
func (b *Buy) GetComputeUnits() uint32 {
	return uint32(b.ComputeUnits)
}
func (b *Buy) GetSlippage() float64 {
	return b.Slippage
}
func (b *Buy) GetSellFee() float64 {
	return b.SellFee
}
