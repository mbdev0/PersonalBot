package dto

type TradingType string

const (
	AFK TradingType = "AFK"
)

type TradingTask struct {
	Type         TradingType `json:"trading_type"`
	BuyAmount    float64     `json:"buy_amount"`
	BuyFee       float64     `json:"buy_fee"`
	ComputeUnits float64     `json:"compute_units"`
	Slippage     float64     `json:"slippage"`
	Wallet       string      `json:"wallet"`
	Filters      Filters     `json:"filters"`
}

type TradingTaskResponse struct {
	Type         TradingType `json:"trading_type"`
	Id           string      `json:"id"`
	BuyAmount    float64     `json:"buy_amount"`
	BuyFee       float64     `json:"buy_fee"`
	ComputeUnits float64     `json:"compute_units"`
	Slippage     float64     `json:"slippage"`
	Filters      Filters     `json:"filters"`
}

type TradingTaskPatch struct {
	Type         *TradingType `json:"trading_type"`
	BuyAmount    *float64     `json:"buy_amount"`
	BuyFee       *float64     `json:"buy_fee"`
	ComputeUnits *float64     `json:"compute_units"`
	Slippage     *float64     `json:"slippage"`
	Wallet       *string      `json:"wallet"`
	Filters      *Filters     `json:"filters"`
}
