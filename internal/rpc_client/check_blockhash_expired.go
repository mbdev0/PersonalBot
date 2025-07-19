package rpcclient

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

func IsBlockhashExpired(lastValidBlockheight uint64) (bool, error) {
	client := GetClient()
	resp, err := client.GetBlockHeight(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return false, err
	}
	return (resp > lastValidBlockheight-150), nil
}
