package client

import (
	"context"
	"pump_fun/internal/core/models"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go/rpc"
)

func IsBlockhashExpired(lastValidBlockheight uint64, cancellationToken models.CancelToken) (bool, error) {
	client := GetClient()
	resp, err := client.GetBlockHeight(cancellationToken.CancellationContext, rpc.CommitmentFinalized)
	if err != nil {
		return false, err
	}
	return (resp > lastValidBlockheight-150), nil
}

func GetLatestBlockhash(cancellation_token models.CancelToken) (*rpc.GetLatestBlockhashResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(cancellation_token.CancellationContext, timeout)
	defer cancel()

	latestHash, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
	}

	return latestHash, nil

}
