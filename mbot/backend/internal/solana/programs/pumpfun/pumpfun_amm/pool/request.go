package pool

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetPoolDataBytes(ctx context.Context, poolAddress solana.PK, httpClient *rpc.Client) ([]byte, error) {
	resp, err := httpClient.GetAccountInfoWithOpts(ctx, poolAddress, &rpc.GetAccountInfoOpts{Encoding: solana.EncodingBase64, Commitment: rpc.CommitmentProcessed})
	if err != nil {
		return nil, err
	}

	return resp.Value.Data.GetBinary(), nil
}
