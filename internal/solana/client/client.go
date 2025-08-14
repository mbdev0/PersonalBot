package client

import (
	"pump_fun/infrastructure/config"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

var (
	clientOnce                   sync.Once
	rlOnce                       sync.Once
	client                       *rpc.Client
	rateLimitClient              *rpc.Client
	numberOfRequestsPerTimeFrame = 1
)

func GetClient() *rpc.Client {
	clientOnce.Do(func() { client = rpc.New(config.GetConfig().HttpNode) })
	return client
}

func GetRatelimtClient() *rpc.Client {
	rlOnce.Do(func() {
		limitClient := rpc.NewWithLimiter(config.GetConfig().HttpNode, rate.Every(time.Second), numberOfRequestsPerTimeFrame)
		rateLimitClient = rpc.NewWithCustomRPCClient(limitClient)
	})
	return rateLimitClient
}
