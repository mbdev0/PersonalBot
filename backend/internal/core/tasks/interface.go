package tasks

type Task interface {
	Id() int64
	SetState(State)
	Type() string
	State() State
	SetSlippage(slippage float64)
	SetComputeUnit(cu uint32)
}
