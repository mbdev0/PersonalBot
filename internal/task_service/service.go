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

var (
	Tasks map[string]tasks.Task = map[string]tasks.Task{}
	mu    sync.Mutex
)

type TaskService struct {
	Executor *transaction.TransactionExecutor
}

func (ts *TaskService) Create(task tasks.Task) (tasks.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	Tasks[task.Id()] = task
	logger.Information(Tasks)
	return task, nil
}

func (ts *TaskService) GetTaskWith(id string) (tasks.Task, error) {
	//hangs here since we're locking a mutex thats already locked
	mu.Lock()
	defer mu.Unlock()
	task, ok := Tasks[id]
	if !ok {
		return nil, fmt.Errorf("Task not found with the id: " + id)
	}
	return task, nil
}

func (ts *TaskService) GetAllTasks() []tasks.Task {
	mu.Lock()
	defer mu.Unlock()
	allTasks := make([]tasks.Task, 0, len(Tasks))
	for _, val := range Tasks {
		allTasks = append(allTasks, val)
	}
	return allTasks
}

func (ts *TaskService) UpdateTask(id string, newTask models.RequestTask) (*tasks.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	task, ok := Tasks[id]
	if !ok {
		return nil, fmt.Errorf("Task not found with the id: " + id)
	}

	err := task.UpdateTask(newTask)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (ts *TaskService) DeleteTask(id string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	_, ok := Tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: " + id)
	}

	delete(Tasks, id)
	return nil
}

func (ts *TaskService) TransistionTask(id string, newState tasks.TaskState) (err error) {
	// in here we'll manage changing state
	task, ok := Tasks[id]
	if !ok {
		return fmt.Errorf("Task not found with the id: " + id)
	}

	// need to verify if is valid transition
	if !transition.IsAbleToTransitionTo(newState, task) {
		return fmt.Errorf("Task not able to transition to next state: " + newState.ToString())
	}

	task.SetState(tasks.State{TaskState: newState})

	switch newState {
	case tasks.TaskRun:
		ts.RunTask(task)
	case tasks.TaskCancel:
		logger.Information("Cancelled Task")
		task.Cancel()
		task.SetState(tasks.State{TaskState: newState})
	}

	return nil
}

func (ts *TaskService) RunTask(task tasks.Task) {
	transactionImpl := ts.Executor.GetImplementation(task)
	go ts.Executor.Execute(transactionImpl)
}
