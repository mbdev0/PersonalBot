package tasks

type Task interface {
	Id() string
	SetState(State)
	Type() string
	State() State
	SetSlippage(slippage float64)
	SetComputeUnit(cu uint32)
}
