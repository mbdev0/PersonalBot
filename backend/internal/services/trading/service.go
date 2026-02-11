package trading

import (
	"context"
	"fmt"
	"maps"
	"personal_bot/app/iterable"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/subscription_hub/strategy"
	taskservice "personal_bot/internal/services/task_service"
	"slices"
	"sync"
	"time"
)

type Service struct {
	strategy    *Strategy
	tasks       map[int64]strategies.Task
	running     map[int64]context.CancelFunc
	mu          *sync.Mutex
	subhub      *strategy.SubscriptionHub
	repo        *repository.TradingRepository
	iterable    *iterable.Iterable
	taskService *taskservice.TaskService
}

func NewTradingService(strat *Strategy, sh *strategy.SubscriptionHub, tr *repository.TradingRepository, ts *taskservice.TaskService) *Service {

	return &Service{
		strategy:    strat,
		tasks:       map[int64]strategies.Task{},
		running:     map[int64]context.CancelFunc{},
		mu:          &sync.Mutex{},
		subhub:      sh,
		repo:        tr,
		iterable:    iterable.NewIterable(),
		taskService: ts,
	}
}

func (s *Service) Create(st strategies.Task) (task strategies.Task, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if st.StrategyTaskId() == 0 {
		st.SetId(s.iterable.ID())
	}

	if _, ok := s.tasks[st.StrategyTaskId()]; ok {
		return nil, fmt.Errorf("task already exists with id: %d", st.StrategyTaskId())
	}

	//if we get a buy/sell strategy we precreate the task to speed up running
	if st.StrategyType() == strategies.BUY {
		err := s.initBuyStrategy(st)
		if err != nil {
			return nil, err
		}
	} else if st.StrategyType() == strategies.SELL {
		err := s.initSellStrategy(st)
		if err != nil {
			return nil, err
		}
	}

	s.tasks[st.StrategyTaskId()] = st
	return st, nil
}

func (s *Service) initBuyStrategy(st strategies.Task) (err error) {
	buyStrategy, ok := st.(*strategies.Buy)
	if !ok {
		return fmt.Errorf("expected strategy task to be a Buy Strategy - ended up getting a different type")
	}

	bt := s.createBuyTask(*buyStrategy)
	createdTask, err := s.taskService.Create(bt)
	if err != nil {
		return err
	}

	buyStrategy.BuyTaskId = createdTask.Id()
	buyStrategy.PositionId = createdTask.Id()
	return nil
}

func (s *Service) createBuyTask(buyTask strategies.Buy) *tasks.BuyTask {
	bt := tasks.NewBuyTask(buyTask.Wallet, buyTask.Token, []tasks.Option{tasks.WithSlippage(buyTask.Slippage), tasks.WithComputeUnits(uint32(buyTask.ComputeUnits))},
		[]tasks.BuyOption{tasks.WithBuyAmount(buyTask.BuyAmount), tasks.WithBuyFee(buyTask.BuyFee), tasks.WithStrategyId(buyTask.StrategyTaskId())},
	)

	return bt
}

func (s *Service) initSellStrategy(st strategies.Task) (err error) {
	sellStrategy, ok := st.(*strategies.Sell)
	if !ok {
		return fmt.Errorf("expected strategy task to be a Sell Strategy - ended up getting a different type")
	}

	sellTask := s.createSellTask(*sellStrategy)
	createdTask, err := s.taskService.Create(sellTask)
	if err != nil {
		return err
	}

	sellStrategy.SellTaskId = createdTask.Id()
	return nil
}

func (s *Service) createSellTask(sell strategies.Sell) *tasks.SellTask {
	st := tasks.NewSellTask(
		sell.GetWallet(),
		sell.Token,
		[]tasks.Option{
			tasks.WithComputeUnits(uint32(sell.GetComputeUnits())),
			tasks.WithSlippage(sell.GetSlippage()),
		},
		[]tasks.SellOption{
			tasks.WithSellAmount(sell.SellAmount),
			tasks.WithSellFee(sell.SellFee),
		},
	)

	return st
}

func (s *Service) Delete(id int64, ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	taskCancel, ok := s.running[id]
	if ok {
		taskCancel()
	}

	task := s.tasks[id]
	delete(s.tasks, id)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	success, err := s.repo.Delete(ctx, id)
	if err != nil {
		s.tasks[id] = task
		return fmt.Errorf("failed to delete from database: %w", err)
	}

	if !success {
		s.tasks[id] = task
		return fmt.Errorf(" unsuccessful delete from database")
	}

	if task.StrategyType() == strategies.BUY {
		task, err := s.taskService.GetTaskWithStrategyId(id)
		if err != nil {
			return err
		}

		err = s.taskService.DeleteTask(task.Id())
		if err != nil {
			return err
		}
	}

	if task.StrategyType() == strategies.SELL {
		sell, ok := task.(*strategies.Sell)
		if !ok {
			return fmt.Errorf("unable to cast task id %d, to type sell strategy", task.StrategyTaskId())
		}

		err := s.taskService.DeleteTask(sell.SellTaskId)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *Service) GetBy(id int64) (strategies.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found with id: %d", id)
	}

	return task, nil
}

func (s *Service) GetAll() []strategies.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	allTasks := make([]strategies.Task, 0, len(s.tasks))
	for _, val := range s.tasks {
		allTasks = append(allTasks, val)
	}

	return allTasks
}

func (s *Service) Update(task strategies.Task, patch strategies.Patch) (strategies.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := patch.ApplyTo(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) Start(id int64) error {
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	_, isRunning := s.running[id]
	if isRunning {
		return fmt.Errorf("task is already running %d", id)
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	task.SetStrategyState(string(strategies.RUNNING))

	switch tsk := task.(type) {
	case *strategies.Afk:
		go s.strategy.AfkSniping(tsk, ctxCancel)
	case *strategies.Buy:
		go s.strategy.Buy(tsk, ctxCancel)
	case *strategies.Sell:
		go s.strategy.Sell(tsk, ctxCancel)
	default:
		//if the task matches no type
		cancel()
		return fmt.Errorf("task doesn't belong to a strategy")
	}

	s.running[task.StrategyTaskId()] = cancel

	return nil
}

func (s *Service) Stop(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// this will call a tasks cancel
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	taskCancel, ok := s.running[id]
	if !ok {
		return fmt.Errorf("task not running with id: %d", id)
	}

	task.SetStrategyState(string(strategies.CREATED))

	taskCancel()
	delete(s.running, task.StrategyTaskId())
	return nil
}

func (s *Service) Subscribe(taskId int64) (*strategy.Subscription, error) {
	sub, err := s.subhub.Subscribe(taskId)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *Service) Unsubscribe(taskId int64) error {
	err := s.subhub.Unsubscribe(taskId)
	return err
}

// loads tasks into memory
func (s *Service) LoadFromDB(ctx context.Context) error {
	ttFromDb, err := s.repo.GetAllTasks(ctx)
	if err != nil {
		return err
	}

	maxId := s.repo.GetMaxId(ctx)
	s.iterable.SetIterable(maxId)

	for _, tt := range ttFromDb {
		s.tasks[tt.StrategyTaskId()] = tt
	}

	return nil
}

// pushes all in-memory tasks to db
func (s *Service) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	deleteSuccess, err := s.repo.DeleteAll(ctx)
	if err != nil {
		return err
	}

	if !deleteSuccess {
		return fmt.Errorf("failed to wipe table for insertion")
	}

	tasksToSave := slices.Collect(maps.Values(s.tasks))

	success, err := s.repo.AddAllTasks(tasksToSave, ctx)
	if err != nil {
		return fmt.Errorf("error whilst shutting down: %w", err)
	}

	if !success {
		return fmt.Errorf("unsuccesful graceful shutdown in tasks service")
	}

	return nil
}
