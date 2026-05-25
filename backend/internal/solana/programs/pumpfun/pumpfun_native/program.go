package pumpfunnative

import (
	"context"
	"math/big"
	"personal_bot/internal/core/position"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	datastream "personal_bot/internal/solana/monitoring/data_stream"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/monitoring/models"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/monitoring"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction/buy"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction/sell"
	"personal_bot/internal/solana/transaction"
)

type PumpFunNative struct {
	datastream datastream.DataStream
}

func New(datastream datastream.DataStream) *PumpFunNative {
	return &PumpFunNative{
		datastream: datastream,
	}
}

func (pn *PumpFunNative) NewBuyTransaction(task *tasks.BuyTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction {
	return buy.NewTransaction(task, posService, publisher)
}

func (pn *PumpFunNative) NewSellTransaction(task *tasks.SellTask, posService position.PositionService, publisher subscriptionhub.Publisher) transaction.Transaction {
	return sell.NewTransaction(task, posService, publisher)
}

func (pn *PumpFunNative) NewMarketCapStream(ctx context.Context, monitoringAddress string, node rpcgroups.RPCNode) <-chan big.Float {
	stream := monitoring.NewPumpfunNativeMarketcap(monitoringAddress, node.WS, node.Http, pn.datastream)
	return stream.StreamMarketCap(ctx)
}

func (pn *PumpFunNative) NewCoinStream(ctx context.Context, node rpcgroups.RPCNode, filters filters.FilterPipeline) <-chan models.Coin {
	stream := monitoring.NewPumpfunNativeCoinCreation(node.WS, node.Http, filters, pn.datastream)
	return stream.StreamCoinCreations(ctx)
}
