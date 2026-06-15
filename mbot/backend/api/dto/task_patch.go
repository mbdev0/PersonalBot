package dto

import programnames "personal_bot/backend/app/program_names"

type PatchRequestTask struct {
	Program           *programnames.Program `json:"program"`
	Type              TransactionType
	Slippage          *float64 `json:"slippage"`
	ComputeUnits      *uint32  `json:"compute_units"`
	WalletAddressName *string  `json:"wallet_name"`
	TokenAddress      *string  `json:"token_address"`
	BuyAmount         *float64 `json:"buy_amount,omitempty"`
	BuyFee            *float64 `json:"buy_fee,omitempty"`
	SellAmount        *float64 `json:"sell_amount,omitempty"`
	SellFee           *float64 `json:"sell_fee,omitempty"`
}
