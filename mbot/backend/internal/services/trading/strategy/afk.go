package strategy

import (
	"context"
	"fmt"
	"personal_bot/backend/internal/core/program"
	rpcgroupsModel "personal_bot/backend/internal/core/rpc_groups"
	"personal_bot/backend/internal/core/strategies"
	"personal_bot/backend/internal/core/tasks"
	"personal_bot/backend/internal/solana/monitoring/filters"
	"personal_bot/backend/internal/solana/monitoring/models"
	"personal_bot/backend/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type AFK struct {
	Strategy
}

func NewAFKEngine(strategy Strategy) *AFK {
	return &AFK{strategy}
}

func (afk *AFK) Run(ctx context.Context, afkTask *strategies.Afk) {
	coins := make(chan models.Coin, 100)
	defer close(coins)

	filterPipeline := filters.FilterPipeline{}
	for _, f := range afkTask.Filters {
		filterPipeline.AddFilter(f())
	}

	node, err := afk.rpcService.GetNode(afkTask.RPCGroupId())
	if err != nil {
		logger.Error(err)
		afkTask.Logger().Error(err)
		return
	}

	prog := program.Resolve(afkTask.Program)
	coinStream := prog.NewCoinStream(ctx, rpcgroupsModel.RPCNode{
		Http: rpc.New(node.Http),
		WS:   node.WS,
	}, filterPipeline)

	logger.Information("started afk monitor")
	afkTask.Logger().Information("started afk monitor")

	counter := 0
	for {
		if counter > 3 {
			break
		}

		select {
		case <-ctx.Done():
			afkTask.SetStrategyState(string(strategies.CANCELLED))
			afk.strategyHub.PublishStateUpdate(afkTask.StrategyTaskId(), afkTask.StrategyState())

			return
		case coin, ok := <-coinStream:
			if !ok {
				afkTask.SetStrategyState(string(strategies.CANCELLED))
				afk.strategyHub.PublishStateUpdate(afkTask.StrategyTaskId(), afkTask.StrategyState())
				return
			}
			logger.Information("found new coin: ", coin.CoinData.Name)
			afkTask.Logger().Information("found new coin: ", coin.CoinData.Name)
			afk.handleNewCoin(ctx, afkTask, coin)
			counter++
		}
	}
}

func (afk *AFK) handleNewCoin(ctx context.Context, afkTask *strategies.Afk, coin models.Coin) {
	bt, err := afk.createAndRunBuyTask(ctx, coin, afkTask)
	if err != nil {
		logger.Error(err)
		afkTask.Logger().Error(err)
		return
	}

	if len(afkTask.SellStrategies) != 0 {
		pos, ok := afk.positionHub.WaitForCreate(bt.Id())
		if !ok {
			logger.Error("timeout whilst waiting for position to be created: ", bt.Id())
			afkTask.Logger().Error("timeout whilst waiting for position to be created: ", bt.Id())

			return
		}

		sellStrats, err := afk.resolveStrategyConfig(afkTask.SellStrategies, pos)
		if err != nil {
			afkTask.Logger().Error(err)
		}
		node, err := afk.rpcService.GetNode(afkTask.RPCGroupId())
		if err != nil {
			logger.Error(err)
			afkTask.Logger().Error(err)
			return
		}

		err = afk.monitorPositionForSellStrategies(ctx, *pos, afkTask, sellStrats, node)
		if err != nil {
			afkTask.Logger().Error(err)
			return
		}
	}

}

func (afk *AFK) createAndRunBuyTask(ctx context.Context, coin models.Coin, afkTask *strategies.Afk) (bt tasks.Task, err error) {
	coinAddr, err := solana.PublicKeyFromBase58(coin.CoinData.TokenAddr)
	if err != nil {
		logger.Error("couldn't read token address correctly: " + err.Error())
		afkTask.Logger().Error("couldn't read token address correctly: " + err.Error())
	}

	node, err := afk.rpcService.GetNode(afkTask.RPCGroupId())
	if err != nil {
		return nil, err
	}

	_, err = afk.rpcService.Load(ctx, afkTask.RPCGroupId())
	if err != nil {
		return nil, err
	}

	t := afk.createBuyTask(afkTask, coinAddr, node)

	bt, err = afk.taskService.Create(t)
	if err != nil {
		return nil, err
	}

	afk.strategyHub.PublishTaskCreation(afkTask.StrategyTaskId(), bt)
	afk.strategyHub.PublishProgressMessage(afkTask.StrategyTaskId(), fmt.Sprintf("created + running coin for %s", coin.CoinData.Symbol))
	afkTask.Logger().Information(fmt.Sprintf("created + running coin for %s", coin.CoinData.Symbol))

	err = afk.taskService.StartTask(bt.Id())
	if err != nil {
		return nil, err
	}

	return bt, nil
}

func (afk *AFK) createBuyTask(afkTask *strategies.Afk, tokenAddr solana.PublicKey, rpc rpcgroupsModel.GroupItem) *tasks.BuyTask {
	bt := tasks.NewBuyTask(afkTask.Wallet, tokenAddr,
		[]tasks.Option{
			tasks.WithProgram(afkTask.GetProgram()),
			tasks.WithSlippage(afkTask.Slippage),
			tasks.WithComputeUnits(uint32(afkTask.ComputeUnits)),
			tasks.WithStrategyId(afkTask.StrategyTaskId()),
			tasks.WithRPCGroupId(afkTask.RPCGroupId()),
			tasks.WithHttpNode(rpc.Http),
			tasks.WithWS(rpc.WS),
			tasks.WithLogger(afkTask.Logger()),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(afkTask.BuyAmount),
			tasks.WithBuyFee(afkTask.BuyFee)},
	)
	return bt
}
