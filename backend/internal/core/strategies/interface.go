package strategies

type Task interface {
	StrategyTaskId() int64
	SetId(id int64)
	StrategyType() TradingType
	StrategyState() string
	SetStrategyState(st string)
	StrategyMessage() string
	SetStrategyMessage(message string)
	RPCGroupId() int64
	GetProgram() string
}
