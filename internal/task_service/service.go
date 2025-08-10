package taskservice

import (
	"fmt"
	"pump_fun/api/models"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/transaction"
	"pump_fun/internal/transition"
	"pump_fun/pkg/logger"
	"sync"
)

type TaskService struct {
	Executor *transaction.TransactionExecutor
	Tasks    map[string]tasks.Task
	mu       sync.Mutex
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

func (ts *TaskService) UpdateTask(id string, newTask models.RequestTask) (tasks.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, ok := ts.Tasks[id]
	if !ok {
		return nil, fmt.Errorf("Task not found with the id: " + id)
	}

	err := task.UpdateTask(newTask)
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

func (ts *TaskService) TransistionTask(id string, newState tasks.TaskState) (err error) {
	// in here we'll manage changing state
	task, ok := ts.Tasks[id]
	if !ok {
		return fmt.Errorf("Task not found with the id: " + id)
	}

	if transition.IsRetryableState(task.GetState().TaskState) {
		task.SetState(tasks.State{TaskState: tasks.TaskCreate, Error: ""})
	}

	// need to verify if is valid transition
	if !transition.IsAbleToTransitionTo(newState, task) {
		return fmt.Errorf("Task not able to transition to next state: " + newState.ToString())
	}

	task.SetState(tasks.State{TaskState: newState})

	switch newState {
	case tasks.TaskRun:
		task.InitCancelToken()
		ts.RunTask(task)
	case tasks.TaskCancel:
		task.SetState(tasks.State{TaskState: newState})
		task.Cancel()
	}

	return nil
}

func (ts *TaskService) RunTask(task tasks.Task) {
	transactionImpl := ts.Executor.GetImplementation(task)
	go ts.Executor.Execute(transactionImpl)
}
