package tasks

type Task interface {
	Id() int64
	SetState(State)
	Type() string
	State() State
	SetId(id int64)
	SetSlippage(slippage float64)
	SetComputeUnit(cu uint32)
	SetTime(t int64)
}
