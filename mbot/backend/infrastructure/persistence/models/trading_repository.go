package models

type TradingRow struct {
	Id              int    `db:"id"`
	TradingType     string `db:"trading_type"`
	WalletId        string `db:"wallet_id"`
	Slippage        int    `db:"slippage"`
	ComputeUnits    int    `db:"compute_units"`
	Config          string `db:"config"`
	TimeCreatedUnix int64  `db:"time_created"`
	RpcGroupId      int64  `db:"rpc_group_id"`
}

type AfkConfig struct {
	Filters Filters `json:"filters"`
	//Stored as Lamports
	BuyFee         int              `json:"buy_fee"`
	BuyAmount      int              `json:"buy_amount"`
	SellFee        *int             `json:"sell_fee"`
	SellStrategies []SellStrategies `json:"sell_strategies"`
	Retries        *uint16          `json:"retries"`
	RetriesDelayMs *uint32          `json:"retries_delay_ms"`
}

type BuyStrategyConfig struct {
	BuyFee         int              `json:"buy_fee"`
	BuyAmount      int              `json:"buy_amount"`
	SellFee        *int             `json:"sell_fee"`
	SellStrategies []SellStrategies `json:"sell_strategies"`
	Token          string           `json:"token_address"`
	BuyTaskId      int              `json:"buy_task_id"`
	PositionId     int              `json:"position_id"`
	Retries        *uint16          `json:"retries"`
	RetriesDelayMs *uint32          `json:"retries_delay_ms"`
}

type SellStrategyConfig struct {
	SellFee        int     `json:"sell_fee"`
	SellAmount     float64 `json:"sell_amount"`
	Token          string  `json:"token_address"`
	SellTaskId     int     `json:"sell_task_id"`
	Retries        *uint16 `json:"retries"`
	RetriesDelayMs *uint32 `json:"retries_delay_ms"`
}

type SpamStrategyConfig struct {
	BuyFee         int     `json:"buy_fee"`
	BuyAmount      int     `json:"buy_amount"`
	SellFee        *int    `json:"sell_fee"`
	Token          string  `json:"token_address"`
	Retries        *uint16 `json:"retries"`
	RetriesDelayMs *uint32 `json:"retries_delay_ms"`
	//formatted as unix timestamp
	StartTime        uint64 `json:"start_time"`
	NumberOfSubTasks uint32 `json:"number_of_sub_tasks"`
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
