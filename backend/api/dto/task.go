package dto

type TransactionType string

const (
	Buy  TransactionType = "Buy"
	Sell TransactionType = "Sell"
)

type RequestTransitionTask struct {
	State string `json:"state"`
}

type RequestTask struct {
	Type              TransactionType `json:"type"`
	Slippage          float64         `json:"slippage"`
	ComputeUnits      uint32          `json:"compute_units"`
	WalletAddressName string          `json:"wallet_name"`
	TokenAddress      string          `json:"token_address"`
	BuyAmount         *float64        `json:"buy_amount,omitempty"`
	BuyFee            *float64        `json:"buy_fee,omitempty"`
	SellAmount        *float64        `json:"sell_amount,omitempty"`
	SellFee           *float64        `json:"sell_fee,omitempty"`
	SellPositionId    *int64          `json:"sell_position_id,omitempty"`
}

type ResponseTask struct {
	TaskId            int64    `json:"task_id"`
	Type              string   `json:"type"`
	Slippage          float64  `json:"slippage"`
	ComputeUnits      uint32   `json:"compute_units"`
	WalletAddressName string   `json:"wallet_name"`
	WalletPublicKey   string   `json:"wallet_address"`
	TokenAddress      string   `json:"token_address"`
	BuyAmount         *float64 `json:"buy_amount,omitempty"`
	BuyFee            *float64 `json:"buy_fee,omitempty"`
	SellAmount        *float64 `json:"sell_amount,omitempty"`
	SellFee           *float64 `json:"sell_fee,omitempty"`
	SellPositionId    *int64   `json:"sell_position_id,omitempty"`
	State             State    `json:"state"`
}

type State struct {
	TaskState string `json:"task_state"`
	Error     string `json:"error"`
}
