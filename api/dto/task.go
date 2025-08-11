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
	Type         TransactionType `json:"type"`
	Slippage     float64         `json:"slippage"`
	ComputeUnits uint32          `json:"compute_units"`
	//this is a security risk -> figure out a way to have these stored in a secure place and loaded from a secure place
	WalletAddressPrivateKey string   `json:"wallet_address_private_key"`
	TokenAddress            string   `json:"token_address"`
	BuyAmount               *float64 `json:"buy_amount,omitempty"`
	BuyFee                  *float64 `json:"buy_fee,omitempty"`
	SellAmount              *float64 `json:"sell_amount,omitempty"`
	SellFee                 *float64 `json:"sell_fee,omitempty"`
}

type ResponseTask struct {
	TaskId       string  `json:"task_id"`
	Type         string  `json:"type"`
	Slippage     float64 `json:"slippage"`
	ComputeUnits uint32  `json:"compute_units"`
	// //this is a security risk -> figure out a way to have these stored in a secure place and loaded from a secure place
	// WalletAddressPrivateKey string   `json:"wallet_address_private_key"`
	TokenAddress string   `json:"token_address"`
	BuyAmount    *float64 `json:"buy_amount,omitempty"`
	BuyFee       *float64 `json:"buy_fee,omitempty"`
	SellAmount   *float64 `json:"sell_amount,omitempty"`
	SellFee      *float64 `json:"sell_fee,omitempty"`
	State        State    `json:"state"`
}

type State struct {
	TaskState string
	Error     string
}
