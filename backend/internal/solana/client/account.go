package client

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var timeout = 10 * time.Second

func GetAccountInfo(address string, ctx context.Context, httpClient *rpc.Client) (*rpc.GetAccountInfoResult, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	accountInfo, err := httpClient.GetAccountInfoWithOpts(ctx, solana.MustPublicKeyFromBase58(address), &rpc.GetAccountInfoOpts{
		Encoding: solana.EncodingBase64, Commitment: rpc.CommitmentProcessed})

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}
