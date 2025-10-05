package dto

type SellStrategyDTO struct {
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	SellAmount float64 `json:"sell_amount"`
}
