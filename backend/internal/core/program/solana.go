package program

import (
	"personal_bot/internal/solana/programs/monitoring"
	"personal_bot/internal/solana/transaction"
)

type SolanaDeps struct {
	BuyTransaction    transaction.Transaction
	SellTransaction   transaction.Transaction
	MarketcapStreamer monitoring.MarketCapMonitor
	NewCoinsStream    monitoring.CoinMonitor
}
