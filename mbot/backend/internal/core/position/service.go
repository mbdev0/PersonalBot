package position

import (
	"context"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type PositionService interface {
	FindPositionIfExists(token solana.PublicKey, walletAddress solana.PublicKey) (*Position, bool)
	GetById(id int64) (*Position, error)
	ReportSell(ctx context.Context, buyTaskId int64, tokensSold *big.Float, solRecieved *big.Float, marketCapSale *big.Float) error
	ReportBuy(ctx context.Context, buyReportPayload ReportBuyPayload) error
}
