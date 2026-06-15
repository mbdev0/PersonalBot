package strategy

import (
	"context"
	"personal_bot/backend/internal/core/strategies"
	"personal_bot/backend/internal/core/tasks"
	"personal_bot/backend/pkg/logger"
	"sync"
	"time"
)

type Spam struct {
	Strategy
}

func NewSpamEngine(strategy Strategy) *Spam {
	return &Spam{strategy}
}

func (s *Spam) Run(ctx context.Context, task *strategies.Spam) {
	select {
	case <-ctx.Done():
		return
	default:
		if task.StartTime == 0 {
			s.noTimerWork(ctx, task)
		} else {
			s.timerWork(ctx, task)
		}
		return
	}

}

func (s *Spam) timerWork(ctx context.Context, task *strategies.Spam) {
	precreatedTasks := []tasks.Task{}

	for i := 0; i < int(task.NumberOfSubTasks); i++ {
		bt, err := s.createBuyTask(task)
		if err != nil {
			logger.Error(err)
			continue
		}
		precreatedTasks = append(precreatedTasks, bt)
	}

	select {
	case <-time.After(time.Until(time.Unix(int64(task.StartTime), 0))):
	case <-ctx.Done():
		s.cancelSubTasks(task)
		return
	}

	for _, t := range precreatedTasks {
		go s.taskService.StartTask(t.Id())
	}

	go func() {
		<-ctx.Done()
		s.cancelSubTasks(task)
	}()
}

func (s *Spam) noTimerWork(ctx context.Context, task *strategies.Spam) {
	var wg sync.WaitGroup
	for i := 0; i < int(task.NumberOfSubTasks); i++ {
		wg.Go(func() {
			s.work(task)
		})
	}

	go func() {
		wg.Wait()
		<-ctx.Done()
		s.cancelSubTasks(task)
	}()
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
		if err != nil {
			logger.Error(err)
		}
	}
}

func (s *Spam) createBuyTask(task *strategies.Spam) (*tasks.BuyTask, error) {
	node, err := s.rpcService.GetNode(task.RPCGroupId())
	if err != nil {
		return nil, err
	}

	opts := []tasks.Option{
		tasks.WithComputeUnits(uint32(task.ComputeUnits)),
		tasks.WithProgram(task.Program),
		tasks.WithRPCGroupId(task.RPCGroupId()),
		tasks.WithHttpNode(node.Http),
		tasks.WithWS(node.WS),
		tasks.WithSlippage(task.Slippage),
		tasks.WithStrategyId(task.StrategyTaskId())}

	if task.Retries != nil {
		opts = append(opts, tasks.WithRetries(*task.Retries))
	}

	if task.RetriesDelayMS != nil {
		opts = append(opts, tasks.WithRetriesDelayMs(*task.RetriesDelayMS))
	}

	bt := tasks.NewBuyTask(
		task.Wallet,
		task.Token,
		opts,
		[]tasks.BuyOption{
			tasks.WithBuyAmount(task.BuyAmount),
			tasks.WithBuyFee(task.BuyFee),
		},
	)

	_, err = s.taskService.Create(bt)
	if err != nil {
		return nil, err
	}

	return bt, nil
}
