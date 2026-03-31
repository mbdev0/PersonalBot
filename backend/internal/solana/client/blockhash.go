package client

import (
	"context"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go/rpc"
)

func IsBlockhashExpired(lastValidBlockheight uint64, rpcClient *rpc.Client, ctx context.Context) (bool, error) {
	resp, err := rpcClient.GetBlockHeight(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return false, err
	}
	return resp > lastValidBlockheight-150, nil
}

func GetLatestBlockhash(rpcClient *rpc.Client, ctx context.Context) (*rpc.GetLatestBlockhashResult, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	latestHash, err := rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return nil, err
	}

	return latestHash, nil

}
