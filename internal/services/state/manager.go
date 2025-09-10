package state

import (
	"context"
	"fmt"
	"pump_fun/internal/core/tasks"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/internal/services/transaction"
	"pump_fun/pkg/logger"
	"sync"
)

type Manager struct {
	executor *transaction.Executor
	running  map[string]context.CancelFunc
	mu       *sync.Mutex
}

func (m *Manager) New(subhub *subscriptionhub.Hub, executor *transaction.Executor) {
	m.running = map[string]context.CancelFunc{}
	m.executor = executor
	m.mu = &sync.Mutex{}
}

func (m *Manager) ExecuteAction(state tasks.TaskState, task tasks.Task) error {
	switch state {
	case tasks.TaskRun:
		if err := m.run(task); err != nil {
			return err
		}
	case tasks.TaskCancel:
		if err := m.cancel(task); err != nil {
			return err
		}
	default:
		return fmt.Errorf("could not find action for task: %s with state: %s", task.Id(), task.State().TaskState.ToString())
	}
	return nil
}

func (m *Manager) run(task tasks.Task) error {
	m.mu.Lock()

	if _, ok := m.running[task.Id()]; ok {
		m.mu.Unlock()
		return fmt.Errorf("task is already running")
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	m.running[task.Id()] = cancel

	m.mu.Unlock()

	transactionImpl, err := m.executor.GetImplementation(task)
	if err != nil {
		return err
	}

	go func() {

		logger.Information(m.running)
		done := make(chan struct{})
		m.executor.Execute(done, transactionImpl, cancelCtx)

		//wait for the channel to close to continue
		<-done

		m.mu.Lock()
		delete(m.running, task.Id())
		logger.Information(m.running)
		m.mu.Unlock()
	}()
	return nil
}

func (m *Manager) cancel(task tasks.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cancel, ok := m.running[task.Id()]

	if !ok {
		return fmt.Errorf("there is no task running with id: %s", task.Id())
	}

	cancel()
	return nil
}
