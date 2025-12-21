package trading

import (
	"context"
	"fmt"
	"pump_fun/internal/core/strategies"
	"sync"
)

type Service struct {
	strategy *Strategy
	tasks    map[int64]strategies.Task
	running  map[int64]context.CancelFunc
	mu       *sync.Mutex
}

func (s *Service) NewTradingService(strat *Strategy) {
	s.strategy = strat
	s.tasks = map[int64]strategies.Task{}
	s.running = map[int64]context.CancelFunc{}
	s.mu = &sync.Mutex{}
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
