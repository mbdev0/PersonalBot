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
	Message      string
	SellTaskId   int64
	TimeCreated  int64
}

func (s *Sell) New() {
	s.strategyType = SELL
}

func (s *Sell) StrategyTaskId() int64 {
	return s.id
}

func (s *Sell) SetId(id int64) {
	s.id = id
}
func (s *Sell) StrategyType() TradingType {
	return s.strategyType
}

func (s *Sell) StrategyState() string {
	return s.State
}

func (s *Sell) SetStrategyState(st string) {
	s.State = st
}

func (s *Sell) StrategyMessage() string {
	return s.Message
}

func (s *Sell) SetStrategyMessage(message string) {
	s.Message = message
}

func (s *Sell) GetWallet() wallets.SolanaWallet {
	return s.Wallet
}
func (s *Sell) GetComputeUnits() uint32 {
	return uint32(s.ComputeUnits)
}
func (s *Sell) GetSlippage() float64 {
	return s.Slippage
}
func (s *Sell) GetSellFee() float64 {
	return s.SellFee
}
