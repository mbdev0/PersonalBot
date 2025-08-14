package tasks

import (
	"context"
)

type Task interface {
	Id() string
	SetState(State)
	Type() string
	State() State
	SetSlippage(slippage float64)
	SetComputeUnit(cu uint32)
	Ctx() context.Context
	Cancel()
	ResetCtx()
}
