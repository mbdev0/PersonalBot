package solanaclient

import (
	"pump_fun/internal/config"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

var httpNodeEndpoint string = config.GetConfig().HttpNode

func NewHttpClient() *rpc.Client {
	rpcClient := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		httpNodeEndpoint,
		rate.Every(time.Second),
		5,
	))

	return rpcClient
}
