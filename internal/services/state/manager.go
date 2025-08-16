package state

import (
	"fmt"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/services/transaction"
	"pump_fun/pkg/logger"
	"sync"
)

type Manager struct {
	executor *transaction.Executor
	running  map[string]tasks.Task
	mu       sync.Mutex
}

func (m *Manager) New() {
	m.running = map[string]tasks.Task{}
	m.executor = &transaction.Executor{}
	m.mu = sync.Mutex{}
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

	m.running[task.Id()] = task
	m.mu.Unlock()

	task.ResetCtx()

	transactionImpl, err := m.executor.GetImplementation(task)
	if err != nil {
		return err
	}

	go func() {
		logger.Information(m.running)
		done := make(chan bool, 1)
		m.executor.Execute(done, transactionImpl)

		isDone := <-done
		if isDone {
			m.mu.Lock()
			delete(m.running, task.Id())
			logger.Information(m.running)
			m.mu.Unlock()
		}
	}()
	return nil
}

func (m *Manager) cancel(task tasks.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Information("in cancel")

	if _, ok := m.running[task.Id()]; !ok {
		return fmt.Errorf("there is no task running with id: %s", task.Id())
	}
	task.Cancel()

	return nil
}
