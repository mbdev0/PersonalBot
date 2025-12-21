package strategies

type Task interface {
	StrategyTaskId() int64
	StrategyType() TradingType
}
