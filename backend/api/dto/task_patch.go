package dto

type PatchRequestTask struct {
	Type         TransactionType
	Slippage     *float64 `json:"slippage"`
	ComputeUnits *uint32  `json:"compute_units"`
	//this is a security risk -> figure out a way to have these stored in a secure place and loaded from a secure place
	WalletAddressName *string  `json:"wallet_name"`
	TokenAddress      *string  `json:"token_address"`
	BuyAmount         *float64 `json:"buy_amount,omitempty"`
	BuyFee            *float64 `json:"buy_fee,omitempty"`
	SellAmount        *float64 `json:"sell_amount,omitempty"`
	SellFee           *float64 `json:"sell_fee,omitempty"`
}
