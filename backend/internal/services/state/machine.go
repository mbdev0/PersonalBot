package state

import (
	"fmt"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/state/transition"
)

type Machine struct {
}

func (m *Machine) Transition(task tasks.Task, newState tasks.TaskState) error {
	if transition.IsRetryableState(task.State().TaskState) {
		task.SetState(tasks.State{TaskState: tasks.TaskCreate, Error: ""})
	}

	if !transition.IsAbleToTransitionTo(newState, task) {
		return fmt.Errorf("Task not able to transition to next state: %s", newState.ToString())
	}

	if newState == tasks.TaskCancel {
		task.SetState(tasks.State{TaskState: tasks.TaskCancel, Error: ""})
		return nil
	}

	err := transition.AutoTransitionTask(task, nil)
	if err != nil {
		return err
	}
	return nil
}
