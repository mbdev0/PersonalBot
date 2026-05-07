package taskservice

import (
	"context"
	"fmt"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/services/transaction"
	"sync"
)

type Manager struct {
	executor *transaction.Executor
	running  map[int64]context.CancelFunc
	mu       *sync.Mutex
}

func NewTaskManager(subhub *subscriptionhub.Hub, executor *transaction.Executor) *Manager {
	return &Manager{
		executor: executor,
		running:  map[int64]context.CancelFunc{},
		mu:       &sync.Mutex{},
	}
}

func (m *Manager) RunTask(task tasks.Task) error {
	m.mu.Lock()

	if _, ok := m.running[task.Id()]; ok {
		m.mu.Unlock()
		return fmt.Errorf("task is already running")
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	m.running[task.Id()] = cancel

	m.mu.Unlock()

	transactionImpl, err := m.executor.GetImplementation(cancelCtx, task)
	if err != nil {
		return err
	}

	go func() {

		done := make(chan struct{})
		m.executor.Execute(cancelCtx, done, transactionImpl)

		//wait for the channel to close to continue
		<-done

		m.mu.Lock()
		delete(m.running, task.Id())
		m.mu.Unlock()
	}()
	return nil
}

func (m *Manager) StopTask(task tasks.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cancel, ok := m.running[task.Id()]

	if !ok {
		return fmt.Errorf("there is no task running with id: %d", task.Id())
	}

	cancel()
	return nil
}
