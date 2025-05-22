package rpcclient

import (
	"context"
	"pump_fun/internal/logger"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetLatestBlockhash() (*rpc.GetLatestBlockhashResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()

	latestHash, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
	}

	return latestHash, nil

}
