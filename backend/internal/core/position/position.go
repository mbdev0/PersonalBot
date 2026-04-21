package position

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type Position struct {
	PositionId           int64
	StrategyId           *int64
	TokenAddress         solana.PublicKey
	WalletAddress        solana.PublicKey
	InitialTokenAmount   *big.Float
	TokenRemaining       *big.Float
	RemainingCostBasis   *big.Float
	FinalizedProfit      *big.Float
	InitialSolanaAmount  *big.Float
	EntryPrice           *big.Float
	MarketCapEntry       *big.Float
	AverageMarketCapExit *big.Float
	TotalMarketCapExit   *big.Float
	NumberOfSells        int
	AddressForUrl        string
}

// we keep floats as strings to preserve some sort of accuracy
type PositionMessage struct {
	MessageType         MessageType `json:"message_type"`
	BuyTaskId           int64       `json:"buy_task_id"`
	StrategyId          *int64      `json:"strategy_id,omitempty"`
	UnrealizedProfit    *big.Float  `json:"unrealized_profit"`
	RealizedProfit      *big.Float  `json:"realized_profit"`
	TotalPnL            *big.Float  `json:"total_pnl"`
	MarketCap           *big.Float  `json:"market_cap"`
	CurrentPrice        *big.Float  `json:"current_price"`
	InitialSolanaAmount *big.Float  `json:"initial_solana_amount"`
	RemainingTokens     *big.Float  `json:"remaining_tokens"`
	Message             string      `json:"message"`
}

type MessageType string

const (
	Created MessageType = "PositionCreated"
	Stopped MessageType = "PositionStopped"
	Update  MessageType = "PositionUpdate"
)

type ReportBuyPayload struct {
	StrategyId    *int64
	BuyTaskId     int64
	TokenAddress  solana.PublicKey
	WalletAddress solana.PublicKey
	TokenAmount   *big.Float
	SolSpent      *big.Float
	MarketCap     *big.Float
	AddressForUrl string
}
