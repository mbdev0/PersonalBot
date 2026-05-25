package program

import (
	"context"
	"math/big"
	"personal_bot/internal/core/position"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/monitoring/models"
	"personal_bot/internal/solana/transaction"
)

// TODO: add unsubsribe methods to marketcapstream + coinstream to kill of goroutines
type Program interface {
	NewBuyTransaction(task *tasks.BuyTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction
	NewSellTransaction(task *tasks.SellTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction
	NewMarketCapStream(ctx context.Context, monitoringAddress string, node rpcgroups.RPCNode) <-chan big.Float
	NewCoinStream(ctx context.Context, node rpcgroups.RPCNode, filters filters.FilterPipeline) <-chan models.Coin
}
