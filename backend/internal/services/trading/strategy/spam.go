package strategy

import (
	"context"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/tasks"
	"personal_bot/pkg/logger"
)

type Spam struct {
	Strategy
}

func NewSpamEngine(strategy Strategy) *Spam {
	return &Spam{strategy}
}

func (s *Spam) Run(ctx context.Context, task *strategies.Spam) {
	for {
		select {
		case <-ctx.Done():
			s.cancelSubTasks(task)
			return
		default:
			for i := 0; i < int(task.NumberOfSubTasks); i++ {
				go s.work(task)
			}
		}
	}
}

func (s *Spam) work(task *strategies.Spam) {
	bt, err := s.createBuyTask(task)
	if err != nil {
		return
	}

	s.taskService.StartTask(bt.Id())
}

func (s *Spam) cancelSubTasks(task *strategies.Spam) {
	tasks := s.taskService.GetTasksWithStrategyId(task.StrategyTaskId())
	for _, t := range tasks {
		err := s.taskService.StopTask(t.Id())
		logger.Error(err)
	}
}

func (s *Spam) createBuyTask(task *strategies.Spam) (*tasks.BuyTask, error) {
	node, err := s.rpcService.GetNode(task.RPCGroupId())
	if err != nil {
		return nil, err
	}

	bt := tasks.NewBuyTask(
		task.Wallet,
		task.Token,
		[]tasks.Option{
			tasks.WithComputeUnits(uint32(task.ComputeUnits)),
			tasks.WithProgram(task.Program),
			tasks.WithRPCGroupId(task.RPCGroupId()),
			tasks.WithHttpNode(node.Http),
			tasks.WithWS(node.WS),
			tasks.WithSlippage(task.Slippage),
			tasks.WithStrategyId(task.StrategyTaskId()),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(task.BuyAmount),
			tasks.WithBuyFee(task.BuyFee),
		},
	)

	return bt, nil
}
