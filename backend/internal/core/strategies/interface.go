package strategies

type Task interface {
	StrategyTaskId() string
	StrategyType() TradingType
}
