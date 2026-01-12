package models

type TaskRow struct {
	Id           int    `db:"id"`
	TaskType     string `db:"task_type"`
	WalletId     string `db:"wallet_id"`
	Slippage     int    `db:"slippage_percentage"`
	ComputeUnits int    `db:"compute_units"`
	Config       string `db:"config"`
	State        string `db:"state"`
	Token        string `db:"token"`
}

// Stored as lamports
type BuyConfig struct {
	BuyFee    int `json:"buy_fee"`
	BuyAmount int `json:"buy_amount"`
}

// Stored as lamports
type SellConfig struct {
	SellFee    int `json:"sell_fee"`
	SellAmount int `json:"sell_percentage"`
}

type AfkConfig struct {
}
