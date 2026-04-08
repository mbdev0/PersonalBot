package client

import (
	"context"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go/rpc"
)

func IsBlockhashExpired(ctx context.Context, lastValidBlockheight uint64, rpcClient *rpc.Client) (bool, error) {
	resp, err := rpcClient.GetBlockHeight(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return false, err
	}
	return resp > lastValidBlockheight-150, nil
}

func GetLatestBlockhash(ctx context.Context, rpcClient *rpc.Client) (*rpc.GetLatestBlockhashResult, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	latestHash, err := rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return nil, err
	}

	return latestHash, nil

}
