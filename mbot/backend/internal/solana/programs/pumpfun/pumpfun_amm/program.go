package pumpfunamm

import (
	"context"
	"math/big"
	"personal_bot/backend/internal/core/position"
	rpcgroups "personal_bot/backend/internal/core/rpc_groups"
	"personal_bot/backend/internal/core/tasks"
	subscriptionhub "personal_bot/backend/internal/services/subscription_hub"
	datastream "personal_bot/backend/internal/solana/monitoring/data_stream"
	"personal_bot/backend/internal/solana/monitoring/filters"
	"personal_bot/backend/internal/solana/monitoring/models"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/monitoring"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/transaction/buy"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/transaction/sell"
	"personal_bot/backend/internal/solana/transaction"
)

type PumpfunAmm struct {
	datastream datastream.DataStream
}

func New(datastream datastream.DataStream) *PumpfunAmm {
	return &PumpfunAmm{
		datastream: datastream,
	}
}

func (pa *PumpfunAmm) NewBuyTransaction(task *tasks.BuyTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction {
	return buy.NewTransaction(task, posService, publisher)
}

func (pa *PumpfunAmm) NewSellTransaction(task *tasks.SellTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction {
	return sell.NewTransaction(task, posService, publisher)
}

func (pa *PumpfunAmm) NewMarketCapStream(ctx context.Context, monitoringAddress string, node rpcgroups.RPCNode) <-chan big.Float {
	stream := monitoring.NewPumpfunAMMMarketCapMonitor(pa.datastream, monitoringAddress, node.WS, node.Http)
	return stream.StreamMarketCap(ctx)
}

func (pa *PumpfunAmm) NewCoinStream(ctx context.Context, node rpcgroups.RPCNode, filters filters.FilterPipeline) <-chan models.Coin {
	stream := monitoring.NewPumpfunAMMCoinCreationMonitor(pa.datastream, node.WS, node.Http, filters)
	return stream.StreamCoinCreation(ctx)
}
