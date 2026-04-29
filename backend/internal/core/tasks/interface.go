package tasks

type Task interface {
	Id() int64
	GetStrategyId() *int64
	SetState(State)
	SetMessage(string)
	Message() string
	Type() string
	State() State
	GetWallet() string
	GetToken() string
}

type Configurable interface {
	SetStrategyId(id int64)
	SetId(id int64)
	SetSlippage(slippage float64)
	SetComputeUnit(cu uint32)
	SetTime(t int64)
	SetHttpNode(rpcNode string)
	SetWSNode(ws string)
	SetRPCGroupId(id int64)
}

type ConfigurableTask interface {
	Task
	Configurable
}
