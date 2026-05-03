package notifier

type BuyNotifierPayload struct {
	TaskType        string
	TaskId          int64
	StrategyId      *int64
	TokensBought    float64
	AmountPaid      float64
	TxHash          string
	WalletAddress   string
	AddressForAxiom string
	TokenAddress    string
}

type SellNotifierPayload struct {
	TaskType        string
	TaskId          int64
	StrategyId      *int64
	TxHash          string
	TokensSold      float64
	TokensRemaining float64
	CurrentProfit   float64
	AmountSold      float64
	WalletAddress   string
	AddressForAxiom string
	TokenAddress    string
}

type ErrorNotifierPayload struct {
	TaskType        string
	TaskId          int64
	StrategyId      *int64
	Error           string
	TxHash          string
	WalletAddress   string
	TokenAddress    string
	AddressForAxiom string
}
