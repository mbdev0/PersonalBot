package strategy

import (
	"context"
	"personal_bot/backend/internal/core/strategies"
	"personal_bot/backend/pkg/logger"
)

type Buy struct {
	Strategy
}

func NewBuyEngine(strategy Strategy) *Buy {
	return &Buy{strategy}
}

func (b *Buy) Run(ctx context.Context, buyTask *strategies.Buy) {
	err := b.taskService.StartTask(buyTask.BuyTaskId)
	if err != nil {
		buyTask.SetStrategyMessage(err.Error())
		logger.Error(err)
	}

	go b.syncStateAndMessage(ctx, buyTask.BuyTaskId, buyTask)

	if len(buyTask.SellStrategies) != 0 {
		pos, ok := b.positionHub.WaitForCreate(buyTask.PositionId)
		if !ok {
			logger.Error("timeout whilst waiting for position to be created: ", buyTask.PositionId)
			return
		}

		sellStrats := b.resolveStrategyConfig(buyTask.SellStrategies, pos)

		node, err := b.rpcService.GetNode(buyTask.RPCGroupId())
		if err != nil {
			logger.Error(err)
			return
		}

		err = b.monitorPositionForSellStrategies(ctx, *pos, buyTask, sellStrats, node)
		if err != nil {
			return
		}
	}
}
