package rpcclient

import (
	"context"
	"pump_fun/internal/models"
	"pump_fun/pkg/logger"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetLatestBlockhash(cancellation_token models.CancelToken) (*rpc.GetLatestBlockhashResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(cancellation_token.CancellationContext, 10000*time.Millisecond)
	defer cancel()

	latestHash, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
	}

	return latestHash, nil

}
