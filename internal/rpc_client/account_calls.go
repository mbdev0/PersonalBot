package rpcclient

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetAccountInfo(address string) (*rpc.GetAccountInfoResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	accountInfo, err := client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(address))

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}

func GetAccountInfoLimited(address string) (*rpc.GetAccountInfoResult, error) {
	client := GetRatelimtClient()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accountInfo, err := client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(address))

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}
