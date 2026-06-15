package strategies

type StrategyConfig struct {
	Type       SellStrategyType
	Value      float64
	SellAmount float64
}

type SellStrategyType string

const (
	StopLossPrice        SellStrategyType = "stop_loss_price"
	StopLossPercentage   SellStrategyType = "stop_loss_percentage"
	StopLossMarketcap    SellStrategyType = "stop_loss_marketcap"
	TakeProfitPrice      SellStrategyType = "take_profit_price"
	TakeProfitPercentage SellStrategyType = "take_profit_percentage"
	TakeProfitMarketCap  SellStrategyType = "take_profit_marketcap"
)
