package dto

type SellStrategyDTO struct {
	Id         string  `json:"id"`
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	SellAmount float64 `json:"sell_amount"`
}
