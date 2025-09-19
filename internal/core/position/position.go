package position

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type Position struct {
	PositionId         string
	TokenAddress       solana.PublicKey
	WalletAddress      solana.PublicKey
	InitialTokenAmount *big.Float
	TokenRemaining     *big.Float
	RemainingCostBasis *big.Float
	FinalizedProfit    *big.Float
	//unfinalized will be calculated on the spot
	// total pnl will be calculated on the spot
}

// we keep floats as strings to preserve some sort of accuracy
type PositionMessage struct {
	MessageType      MessageType `json:"message_type"`
	BuyTaskId        string      `json:"buy_task_id"`
	UnrealizedProfit string      `json:"unrealized_profit"`
	RealizedProfit   string      `json:"realized_profit"`
	TotalPnL         string      `json:"total_pnl"`
	MarketCap        string      `json:"market_cap"`
	RemainingTokens  string      `json:"remaining_tokens"`
	Message          string      `json:"message"`
}

type MessageType string

const (
	Created MessageType = "PositionCreated"
	Stopped MessageType = "PositionStopped"
	Update  MessageType = "PositionUpdate"
)
