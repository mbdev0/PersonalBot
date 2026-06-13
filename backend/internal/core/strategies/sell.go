package strategies

import (
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go"
)

type Sell struct {
	id             int64
	Program        string
	strategyType   TradingType
	SellAmount     float64
	SellFee        float64
	Slippage       float64
	ComputeUnits   float64
	Token          solana.PublicKey
	Wallet         wallets.SolanaWallet
	State          string
	Message        string
	SellTaskId     int64
	TimeCreated    int64
	RPCGroup       rpcgroups.RPCGroup
	Retries        *uint16
	RetriesDelayMS *uint32
}

func (s *Sell) New() {
	s.strategyType = SELL
}

func (s *Sell) GetProgram() string {
	return s.Program
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

func (s *Sell) RPCGroupId() int64 {
	return s.RPCGroup.Id
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

func (s *Sell) GetTimeCreated() int64 {
	return s.TimeCreated
}

func (s *Sell) SetTimeCreated(time int64) {
	s.TimeCreated = time
}
