package dto

import programnames "personal_bot/app/program_names"

type TransactionType string

const (
	Buy  TransactionType = "Buy"
	Sell TransactionType = "Sell"
)

type ActionType string

const (
	Run  ActionType = "Run"
	Stop ActionType = "Stop"
)

type RequestTransitionTask struct {
	Action ActionType `json:"action"`
}

type RequestTask struct {
	Program           programnames.Program `json:"program"`
	Type              TransactionType      `json:"type"`
	Slippage          float64              `json:"slippage"`
	ComputeUnits      uint32               `json:"compute_units"`
	WalletAddressName string               `json:"wallet_name"`
	TokenAddress      string               `json:"token_address"`
	BuyAmount         *float64             `json:"buy_amount,omitempty"`
	BuyFee            *float64             `json:"buy_fee,omitempty"`
	SellAmount        *float64             `json:"sell_amount,omitempty"`
	SellFee           *float64             `json:"sell_fee,omitempty"`
	SellPositionId    *int64               `json:"sell_position_id,omitempty"`
	StrategyId        *int64               `json:"strategy_id,omitempty"`
	RPCGroupId        int64                `json:"rpc_group_id"`
}

type ResponseTask struct {
	Program           string   `json:"program"`
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
	StrategyId        *int64   `json:"strategy_id,omitempty"`
	State             State    `json:"state"`
	Message           string   `json:"message"`
	TimeCreated       int64    `json:"time_created"`
	HTTPNode          string   `json:"http_node"`
	WSNode            string   `json:"ws_node"`
	RPCGroupId        int64    `json:"rpc_group_id"`
}

type State struct {
	TaskState string `json:"task_state"`
	Error     string `json:"error"`
}
