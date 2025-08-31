package trading

import (
	"context"
	"fmt"
	"pump_fun/internal/core/strategies"
	"sync"
)

type Service struct {
	strategy *Strategy
	tasks    map[string]strategies.Task
	running  map[string]context.CancelFunc
	mu       *sync.Mutex
}

func (s *Service) NewTradingService(strat *Strategy) {
	s.strategy = strat
	s.tasks = map[string]strategies.Task{}
	s.running = map[string]context.CancelFunc{}
	s.mu = &sync.Mutex{}
}

func (s *Service) Create(st strategies.Task) (task strategies.Task, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[st.StrategyTaskId()]; ok {
		return nil, fmt.Errorf("task already exists with id: %s", st.StrategyTaskId())
	}

	s.tasks[st.StrategyTaskId()] = st
	return st, nil
}

func (s *Service) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found with id: %s", id)
	}

	delete(s.tasks, id)
	return nil

}

func (s *Service) GetBy(id string) (strategies.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found with id: %s", id)
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

func (s *Service) Start(id string, ctx context.Context) error {
	// this will probably be a switch case over the task type and run the function underneath
	// we will do a go func with a ctx that we derive from the request
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %s", id)
	}

	ctxCancel, cancel := context.WithCancel(ctx)

	switch tsk := task.(type) {
	case *strategies.Afk:
		tsk.Cancel = cancel
		go s.strategy.AfkSniping(tsk, ctxCancel)
		s.running[task.StrategyTaskId()] = tsk.Cancel
	default:
		//if the task matches no type
		cancel()
		return fmt.Errorf("task doesn't belong to a strategy")
	}

	return nil
}

func (s *Service) Stop(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// this will call a tasks cancel
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %s", id)
	}

	taskCancel, ok := s.running[id]
	if !ok {
		return fmt.Errorf("task not running with id: %s", id)
	}

	taskCancel()
	delete(s.running, task.StrategyTaskId())
	return nil
}
