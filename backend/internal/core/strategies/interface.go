package strategies

type Task interface {
	StrategyTaskId() int64
	SetId(id int64)
	StrategyType() TradingType
}
