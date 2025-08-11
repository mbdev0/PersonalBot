package state

import (
	"fmt"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/services/state/transition"
	"pump_fun/internal/services/transaction"
)

type Machine struct {
	executor transaction.Executor
}

func (m *Machine) NewStateMachine() {
	m.executor = transaction.Executor{}
}

func (m *Machine) Transition(task tasks.Task, newState tasks.TaskState) error {
	if transition.IsRetryableState(task.GetState().TaskState) {
		task.SetState(tasks.State{TaskState: tasks.TaskCreate, Error: ""})
	}

	if !transition.IsAbleToTransitionTo(newState, task) {
		return fmt.Errorf("Task not able to transition to next state: " + newState.ToString())
	}

	switch newState {
	case tasks.TaskRun:
		task.InitCancelToken()
		if err := m.runTask(task); err != nil {
			return err
		}
		task.SetState(tasks.State{TaskState: tasks.TaskRun})
	case tasks.TaskCancel:
		task.Cancel()
		task.SetState(tasks.State{TaskState: tasks.TaskCancel})
	default:
		task.SetState(tasks.State{TaskState: newState})
	}

	return nil
}

func (m *Machine) runTask(task tasks.Task) error {
	transactionImpl, err := m.executor.GetImplementation(task)
	if err != nil {
		return err
	}
	go m.executor.Execute(transactionImpl)
	return nil
}
