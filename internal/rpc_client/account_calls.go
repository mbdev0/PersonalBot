package rpcclient

import (
	"context"
	"pump_fun/internal/core/models"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var timeout = 10 * time.Second

func GetAccountInfo(address string, cancellation_token models.CancelToken) (*rpc.GetAccountInfoResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(cancellation_token.CancellationContext, timeout)
	defer cancel()

	accountInfo, err := client.GetAccountInfoWithOpts(ctx, solana.MustPublicKeyFromBase58(address), &rpc.GetAccountInfoOpts{
		Encoding: solana.EncodingBase64})

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
