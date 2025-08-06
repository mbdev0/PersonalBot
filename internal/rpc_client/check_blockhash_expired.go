package rpcclient

import (
	"pump_fun/internal/models"

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
