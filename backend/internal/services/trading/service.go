package trading

import (
	"context"
	"fmt"
	"maps"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/services/subscription_hub/strategy"
	"slices"
	"sync"
)

type Service struct {
	strategy *Strategy
	tasks    map[int64]strategies.Task
	running  map[int64]context.CancelFunc
	mu       *sync.Mutex
	subhub   *strategy.SubscriptionHub
	repo     *repository.TradingRepository
}

func (s *Service) NewTradingService(strat *Strategy, sh *strategy.SubscriptionHub, tr *repository.TradingRepository) {
	s.strategy = strat
	s.tasks = map[int64]strategies.Task{}
	s.running = map[int64]context.CancelFunc{}
	s.mu = &sync.Mutex{}
	s.subhub = sh
	s.repo = tr
}

func (s *Service) Create(st strategies.Task) (task strategies.Task, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[st.StrategyTaskId()]; ok {
		return nil, fmt.Errorf("task already exists with id: %d", st.StrategyTaskId())
	}

	s.tasks[st.StrategyTaskId()] = st
	return st, nil
}

func (s *Service) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	taskCancel, ok := s.running[id]
	if ok {
		taskCancel()
	}

	delete(s.tasks, id)
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

	switch tsk := task.(type) {
	case *strategies.Afk:
		go s.strategy.AfkSniping(tsk, ctxCancel)
		s.running[task.StrategyTaskId()] = cancel
	default:
		//if the task matches no type
		cancel()
		return fmt.Errorf("task doesn't belong to a strategy")
	}

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
