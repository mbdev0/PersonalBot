package strategies

import (
	"math/big"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go"
)

type Buy struct {
	Program        string
	id             int64
	strategyType   TradingType
	BuyAmount      *big.Int
	BuyFee         float64
	Slippage       float64
	ComputeUnits   float64
	Token          solana.PublicKey
	Wallet         wallets.SolanaWallet
	SellStrategies []StrategyConfig
	SellFee        *float64
	BuyTaskId      int64
	PositionId     int64
	State          string
	Message        string
	TimeCreated    int64
	RPCGroup       rpcgroups.RPCGroup
	Retries        uint16
	RetriesDelayMS uint32
}

func (b *Buy) New() {
	b.strategyType = BUY
}

func (b *Buy) GetProgram() string {
	return b.Program
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

func (b *Buy) StrategyMessage() string {
	return b.Message
}

func (b *Buy) RPCGroupId() int64 {
	return b.RPCGroup.Id
}

func (b *Buy) SetStrategyMessage(message string) {
	b.Message = message
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
func (b *Buy) GetSellFee() *float64 {
	return b.SellFee
}
