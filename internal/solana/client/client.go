package client

import (
	"pump_fun/internal/launch/config"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

var lock = &sync.Mutex{}

var client *rpc.Client

var rateLimitClient *rpc.Client

var numberOfRequestsPerTimeFrame = 1

func GetClient() *rpc.Client {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()

		if client == nil {
			client = rpc.New(config.GetConfig().HttpNode)
		}
	}

	return client
}

func GetRatelimtClient() *rpc.Client {
	if rateLimitClient == nil {
		lock.Lock()
		defer lock.Unlock()

		if rateLimitClient == nil {
			limitClient := rpc.NewWithLimiter(config.GetConfig().HttpNode, rate.Every(time.Second), numberOfRequestsPerTimeFrame)
			rateLimitClient = rpc.NewWithCustomRPCClient(limitClient)
		}
	}

	return rateLimitClient
}
