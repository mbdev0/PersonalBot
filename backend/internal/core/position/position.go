package position

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type Position struct {
	PositionId          int64
	TokenAddress        solana.PublicKey
	WalletAddress       solana.PublicKey
	InitialTokenAmount  *big.Float
	TokenRemaining      *big.Float
	RemainingCostBasis  *big.Float
	FinalizedProfit     *big.Float
	InitialSolanaAmount *big.Float
	EntryPrice          *big.Float
}

// we keep floats as strings to preserve some sort of accuracy
type PositionMessage struct {
	MessageType         MessageType `json:"message_type"`
	BuyTaskId           int64       `json:"buy_task_id"`
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
