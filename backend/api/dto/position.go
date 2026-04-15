package dto

type PositionDto struct {
	PositionId           int64  `json:"position_id"`
	TokenAddress         string `json:"token_address"`
	WalletAddress        string `json:"wallet_address"`
	InitialTokenAmount   string `json:"initial_token_amount"`
	TokenRemaining       string `json:"tokens_remaining"`
	RemaningCostBasis    string `json:"remaining_cost_basis"`
	FinalizedProfit      string `json:"finalized_profit"`
	MarketCapEntry       string `json:"market_cap_entry_usd"`
	AverageMarketCapExit string `json:"average_market_cap_exit_usd"`
}

type PositionDashboard []PositionDashboardRow

type PositionDashboardRow struct {
	TotalPNL              string //we can get
	AverageMarketCapEntry string //we can get 
	AverageMarketCapExit  string //we can get
	Coin                  string //we can get
}
