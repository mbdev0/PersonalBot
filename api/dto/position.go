package dto

type PositionDto struct {
	PositionId         string `json:"position_id"`
	TokenAddress       string `json:"token_address"`
	WalletAddress      string `json:"wallet_address"`
	InitialTokenAmount string `json:"initial_token_amount"`
	TokenRemaining     string `json:"tokens_remaining"`
	RemaningCostBasis  string `json:"remaining_cost_basis"`
	FinalizedProfit    string `json:"finalized_profit"`
}
