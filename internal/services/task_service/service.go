package taskservice

import (
	"fmt"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/services/state"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/pkg/logger"
	"sync"
)

type TaskService struct {
	StateMachine *state.Machine
	StateManager *state.Manager
	Hub          *subscriptionhub.Hub
	Tasks        map[string]tasks.Task
	mu           sync.Mutex
}

func (ts *TaskService) NewTaskService() {
	ts.Tasks = map[string]tasks.Task{}
}

func (ts *TaskService) Create(task tasks.Task) (tasks.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.Tasks[task.Id()] = task
	logger.Information(ts.Tasks)
	return task, nil
}

func (ts *TaskService) GetTaskWith(id string) (tasks.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, ok := ts.Tasks[id]
	if !ok {
		return nil, fmt.Errorf("Task not found with the id: " + id)
	}
	return task, nil
}

func (ts *TaskService) GetAllTasks() []tasks.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	allTasks := make([]tasks.Task, 0, len(ts.Tasks))
	for _, val := range ts.Tasks {
		allTasks = append(allTasks, val)
	}
	return allTasks
}

func (ts *TaskService) UpdateTask(task tasks.Task, patch tasks.TaskPatch) (tasks.Task, error) {
	err := patch.ApplyTo(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) DeleteTask(id string) (err error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	_, ok := ts.Tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: " + id)
	}

	delete(ts.Tasks, id)
	return nil
}

func (ts *TaskService) TransitionTask(id string, newState tasks.TaskState) (err error) {
	// in here we'll manage changing state
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, ok := ts.Tasks[id]
	if !ok {
		return fmt.Errorf("Task not found with the id: " + id)
	}

	err = ts.StateMachine.Transition(task, newState)
	ts.Hub.PublishStateChange(task)

	if err != nil {
		return fmt.Errorf("transition failed for task %s with error: %w", task.Id(), err)
	}

	err = ts.StateManager.ExecuteAction(newState, task)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) Subscribe(task tasks.Task) (*subscriptionhub.Subscription, error) {
	c, err := ts.Hub.Subscribe(task)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (ts *TaskService) Unsubscribe(task tasks.Task) error {
	err := ts.Hub.Unsubcribe(task)
	if err != nil {
		return err
	}
	return nil
}
