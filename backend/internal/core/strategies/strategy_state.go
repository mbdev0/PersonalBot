package strategies

type StrategyState string

const (
	CREATED StrategyState = "CREATED"
	RUNNING StrategyState = "RUNNING"
	FAILED  StrategyState = "FAILED"
	SUCCESS StrategyState = "SUCCESS"
)
