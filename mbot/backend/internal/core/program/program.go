package program

import (
	"context"
	"math/big"
	"personal_bot/backend/internal/core/position"
	rpcgroups "personal_bot/backend/internal/core/rpc_groups"
	"personal_bot/backend/internal/core/tasks"
	subscriptionhub "personal_bot/backend/internal/services/subscription_hub"
	"personal_bot/backend/internal/solana/monitoring/filters"
	"personal_bot/backend/internal/solana/monitoring/models"
	"personal_bot/backend/internal/solana/transaction"
)

type Program interface {
	NewBuyTransaction(task *tasks.BuyTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction
	NewSellTransaction(task *tasks.SellTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction
	NewMarketCapStream(ctx context.Context, monitoringAddress string, node rpcgroups.RPCNode) <-chan big.Float
	NewCoinStream(ctx context.Context, node rpcgroups.RPCNode, filters filters.FilterPipeline) <-chan models.Coin
}
