package taskservice

import (
	"fmt"
	"pump_fun/api/models"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/transaction"
	"pump_fun/pkg/logger"
	"sync"
)

var (
	Tasks map[string]*tasks.Task = map[string]*tasks.Task{}
	mu    sync.Mutex
)

type TaskService struct {
	Executor *transaction.TransactionExecutor
}

func (ts *TaskService) Create(task tasks.Task) (tasks.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	Tasks[task.Id()] = &task
	logger.Information(Tasks)
	return task, nil
}

func (ts *TaskService) GetTaskWith(id string) (*tasks.Task, error) {
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
		allTasks = append(allTasks, *val)
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
	t := *task
	err := t.UpdateTask(newTask)
	if err != nil {
		return nil, err
	}

	return task, nil
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
	taskPtr, ok := Tasks[id]
	if !ok {
		return fmt.Errorf("Task not found with the id: " + id)
	}

	// need to verify if is valid transition

	task := *taskPtr
	task.SetState(newState)

	switch newState {
	case tasks.TaskStateRunning:
		ts.RunTask(task)
	case tasks.TaskStateCancelled:
		logger.Information("Cancelled Task")
		task.SetState(tasks.TaskStateCreated)
	}

	return nil
}

func (ts *TaskService) RunTask(task tasks.Task) (tasks.Task, error) {
	switch task.GetTaskType() {
	case "Buy":
		logger.Information("Buy Task")
	case "Sell":
		logger.Information("SellTask")
	}

	ts.Executor.Execute(task)
	return task, nil
}
