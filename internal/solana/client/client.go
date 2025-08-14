package client

import (
	"pump_fun/infrastructure/config"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

var (
	once                         sync.Once
	client                       *rpc.Client
	rateLimitClient              *rpc.Client
	numberOfRequestsPerTimeFrame = 1
)

func GetClient() *rpc.Client {
	once.Do(func() { client = rpc.New(config.GetConfig().HttpNode) })
	return client
}

func GetRatelimtClient() *rpc.Client {
	once.Do(func() {
		limitClient := rpc.NewWithLimiter(config.GetConfig().HttpNode, rate.Every(time.Second), numberOfRequestsPerTimeFrame)
		rateLimitClient = rpc.NewWithCustomRPCClient(limitClient)
	})
	return rateLimitClient
}
