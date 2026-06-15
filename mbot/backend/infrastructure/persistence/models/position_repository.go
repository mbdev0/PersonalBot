package models

type PositionRepository struct {
	PositionId         int64  `db:"position_id"`
	StrategyId         *int64 `db:"strategy_id"`
	TokenAddress       string `db:"token_address"`
	WalletAddress      string `db:"wallet_address"`
	InitialTokenAmount string `db:"init_token_amount"`
	TokensRemaining    string `db:"token_remaining"`
	RemainingCostBasis string `db:"remaining_cost_basis"`
	FinalizedProfit    string `db:"finalized_profit"`
	InitialAmount      string `db:"initial_amount"`
	EntryPrice         string `db:"entry_price"`
	McapEntry          string `db:"mcap_entry"`
	AverageMcapExit    string `db:"average_mcap_exit"`
	TotalMarketCapExit string `db:"total_marketcap_exit"`
	NumberOfSells      int64  `db:"number_of_sales"`
	AddressForUrl      string `db:"address_for_url"`
}
