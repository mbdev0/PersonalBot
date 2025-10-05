package client

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var timeout = 10 * time.Second

func GetAccountInfo(address string, ctx context.Context) (*rpc.GetAccountInfoResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	accountInfo, err := client.GetAccountInfoWithOpts(ctx, solana.MustPublicKeyFromBase58(address), &rpc.GetAccountInfoOpts{
		Encoding: solana.EncodingBase64, Commitment: rpc.CommitmentProcessed})

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}

func GetAccountInfoLimited(address string) (*rpc.GetAccountInfoResult, error) {
	client := GetRatelimtClient()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	accountInfo, err := client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(address))

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}

// func GetAccountInfoUntilResponse(address string, ctx context.Context) (*rpc.GetAccountInfoResult, error) {
// 	client := GetClient()
// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	for {
// 		accountInfo, err
// 	}

// }
