package dto

type TradingType string

const (
	AFK TradingType = "AFK"
)

type TradingTask struct {
	Type           TradingType       `json:"trading_type"`
	BuyAmount      float64           `json:"buy_amount"`
	BuyFee         float64           `json:"buy_fee"`
	ComputeUnits   float64           `json:"compute_units"`
	Slippage       float64           `json:"slippage"`
	WalletName     string            `json:"wallet_name"`
	Filters        Filters           `json:"filters"`
	SellStrategies []SellStrategyDTO `json:"sell_strategies"`
	SellFee        float64           `json:"sell_fee"`
}

type TradingTaskResponse struct {
	Type           TradingType       `json:"trading_type"`
	Id             int64             `json:"id"`
	BuyAmount      float64           `json:"buy_amount"`
	BuyFee         float64           `json:"buy_fee"`
	ComputeUnits   float64           `json:"compute_units"`
	Slippage       float64           `json:"slippage"`
	WalletName     string            `json:"wallet_name"`
	WalletAddress  string            `json:"wallet_address"`
	Filters        Filters           `json:"filters"`
	SellStrategies []SellStrategyDTO `json:"sell_strategies"`
	SellFee        float64           `json:"sell_fee"`
}

type TradingTaskPatch struct {
	Type           *TradingType       `json:"trading_type"`
	BuyAmount      *float64           `json:"buy_amount"`
	BuyFee         *float64           `json:"buy_fee"`
	ComputeUnits   *float64           `json:"compute_units"`
	Slippage       *float64           `json:"slippage"`
	Wallet         *string            `json:"wallet_name"`
	Filters        *Filters           `json:"filters"`
	SellStrategies *[]SellStrategyDTO `json:"sell_strategies"`
	SellFee        *float64           `json:"sell_fee"`
}
