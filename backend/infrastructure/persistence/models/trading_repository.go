package models

type TradingRow struct {
	Id           int    `db:"id"`
	TradingType  string `db:"trading_type"`
	WalletId     string `db:"wallet_id"`
	Slippage     int    `db:"slippage"`
	ComputeUnits int    `db:"compute_units"`
	Config       string `db:"config"`
}

type AfkConfig struct {
	Filters Filters `json:"filters"`
	//Stored as Lamports
	BuyFee         int              `json:"buy_fee"`
	BuyAmount      int              `json:"buy_amount"`
	SellFee        int              `json:"sell_fee"`
	SellStrategies []SellStrategies `json:"sell_strategies"`
}

type BuyStrategyConfig struct {
	BuyFee         int              `json:"buy_fee"`
	BuyAmount      int              `json:"buy_amount"`
	SellFee        int              `json:"sell_fee"`
	SellStrategies []SellStrategies `json:"sell_strategies"`
	Token          string           `json:"token_address"`
	BuyTaskId      int              `json:"buy_task_id"`
	PositionId     int              `json:"position_id"`
}

type Filters struct {
	HasWebsite  *bool   `json:"has_website"`
	HasTwitter  *bool   `json:"has_twitter"`
	HasTelegram *bool   `json:"has_telegram"`
	DevWallet   *string `json:"dev_wallet"`
}

type SellStrategies struct {
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	SellAmount float64 `json:"sell_amount"`
}
