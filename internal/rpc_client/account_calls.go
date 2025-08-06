package rpcclient

import (
	"context"
	"pump_fun/internal/models"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetAccountInfo(address string, cancellation_token models.CancelToken) (*rpc.GetAccountInfoResult, error) {
	client := GetClient()

	ctx, cancel := context.WithTimeout(cancellation_token.CancellationContext, 10000*time.Millisecond)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accountInfo, err := client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(address))

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}
