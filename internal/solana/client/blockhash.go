package client

import (
	"context"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go/rpc"
)

func IsBlockhashExpired(lastValidBlockheight uint64, ctx context.Context) (bool, error) {
	client := GetClient()
	resp, err := client.GetBlockHeight(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return false, err
	}
	return resp > lastValidBlockheight-150, nil
}

func GetLatestBlockhash(ctx context.Context) (*rpc.GetLatestBlockhashResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	latestHash, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return nil, err
	}

	return latestHash, nil

}
