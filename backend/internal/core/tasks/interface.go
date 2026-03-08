package tasks

import "github.com/gagliardetto/solana-go"

type Task interface {
	Id() int64
	GetStrategyId() *int64
	SetStrategyId(id int64)
	SetState(State)
	SetMessage(string)
	Message() string
	Type() string
	State() State
	SetId(id int64)
	SetSlippage(slippage float64)
	SetComputeUnit(cu uint32)
	SetTime(t int64)
	GetWallet() solana.PublicKey
	GetToken() solana.PublicKey
}
