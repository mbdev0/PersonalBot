package strategies

type TradingType string

const (
	AFK  TradingType = "AFK"
	BUY  TradingType = "BUY"
	SELL TradingType = "SELL"
	SPAM TradingType = "SPAM"
)
