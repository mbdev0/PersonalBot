package solanaclient

import (
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

var httpNodeEndpoint string = rpc.MainNetBeta.RPC 	// change as needed -> temporary solution until we have a config file set up

func NewHttpClient() *rpc.Client {
	rpcClient := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		httpNodeEndpoint,
		rate.Every(time.Second),
		5,
	))

	return rpcClient
}
