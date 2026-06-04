package strategy

import (
	"context"
	"personal_bot/internal/core/strategies"
	"personal_bot/pkg/logger"
)

type Sell struct {
	Strategy
}

func NewSellEngine(strategy Strategy) *Sell {
	return &Sell{strategy}
}

func (s *Sell) Run(ctxCancel context.Context, tsk *strategies.Sell) {
	if err := s.taskService.StartTask(tsk.SellTaskId); err != nil {
		logger.Error(err)
	}
	go s.syncStateAndMessage(ctxCancel, tsk.SellTaskId, tsk)
}
