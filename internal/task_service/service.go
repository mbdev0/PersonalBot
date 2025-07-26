package taskservice

import (
	"fmt"
	"pump_fun/internal/models/tasks"
	"pump_fun/pkg/logger"
	"sync"
)

var (
	Tasks map[string]*tasks.Task = map[string]*tasks.Task{}
	mu    sync.Mutex
)

type TaskService struct {
}

func (ts *TaskService) Create(task tasks.Task) (tasks.Task, error) {
	mu.Lock()
	Tasks[task.Id()] = &task
	mu.Unlock()

	logger.Information(Tasks)
	return task, nil
}

func (ts *TaskService) RunTask(task tasks.Task) (tasks.Task, error) {
	switch task.(type) {
	case *tasks.BuyTask:
		fmt.Println("buy task")
		logger.Information("Buy Task")
	case *tasks.SellTask:
		fmt.Println("sell task")
		logger.Information("SellTask")
	}

	return task, nil
}
