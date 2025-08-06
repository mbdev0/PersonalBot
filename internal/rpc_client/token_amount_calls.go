package rpcclient

import (
	"pump_fun/internal/models"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTokenAccountBalance(associatedTokenAddress solana.PublicKey, cancellationToken models.CancelToken) (tokenAmount *uint64, err error) {
	client := GetClient()

	result, err := client.GetTokenAccountBalance(cancellationToken.CancellationContext, associatedTokenAddress, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseInt(result.Value.Amount, 10, 64)
	if err != nil {
		return nil, err
	}

	uintAmount := uint64(amount)

	return &uintAmount, nil
}
