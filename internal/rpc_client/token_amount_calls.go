package rpcclient

import (
	"context"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTokenAccountBalance(associatedTokenAddress solana.PublicKey) (tokenAmount *uint64, err error) {
	client := GetClient()

	result, err := client.GetTokenAccountBalance(context.Background(), associatedTokenAddress, rpc.CommitmentConfirmed)
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
