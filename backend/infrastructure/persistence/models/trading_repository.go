type TradingRow struct {
	Id           int    `db:"id"`
	TradingType  string `db:"trading_type"`
	WalletId     int    `db:"wallet_id"`
	Slippage     int    `db:"slippage"`
	ComputeUnits int    `db:"compute_units"`
	Config       string `db:"config"`
}

type AfkConfig struct {
	Filters        string `json:"filters"` // TODO - Change this to a struct of filters
	BuyAmount      int    `json:"buy_amount"`
	BuyFee         int    `json:"buy_fee"`
	SellStrategies string `json:"sell_strategies"` //TODO - change this to a struct of sell strats
}