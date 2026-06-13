package strategies

import (
	"math/big"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go"
)

type Spam struct {
	Program          string
	id               int64
	strategyType     TradingType
	BuyAmount        *big.Int
	BuyFee           float64
	Token            solana.PublicKey
	Slippage         float64
	ComputeUnits     float64
	Wallet           wallets.SolanaWallet
	SellFee          *float64
	State            string
	Message          string
	TimeCreated      int64
	RPCGroup         rpcgroups.RPCGroup
	Retries          *uint16
	RetriesDelayMS   *uint32
	NumberOfSubTasks uint32
	//formatted as unix timestamp
	StartTime uint64
}

func (s *Spam) New() {
	s.strategyType = SPAM
}

func (s *Spam) GetProgram() string {
	return s.Program
}

func (s *Spam) StrategyTaskId() int64 {
	return s.id
}

func (s *Spam) SetId(id int64) {
	s.id = id
}

func (s *Spam) StrategyType() TradingType {
	return s.strategyType
}

func (s *Spam) StrategyState() string {
	return s.State
}

func (s *Spam) SetStrategyState(st string) {
	s.State = st
}

func (s *Spam) StrategyMessage() string {
	return s.Message
}

func (s *Spam) RPCGroupId() int64 {
	return s.RPCGroup.Id
}

func (s *Spam) SetStrategyMessage(message string) {
	s.Message = message
}

func (s *Spam) GetWallet() wallets.SolanaWallet {
	return s.Wallet
}
func (s *Spam) GetComputeUnits() uint32 {
	return uint32(s.ComputeUnits)
}
func (s *Spam) GetSlippage() float64 {
	return s.Slippage
}
func (s *Spam) GetSellFee() *float64 {
	return s.SellFee
}

func (s *Spam) GetTimeCreated() int64 {
	return s.TimeCreated
}

func (s *Spam) SetTimeCreated(time int64) {
	s.TimeCreated = time
}
