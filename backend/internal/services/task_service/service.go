package taskservice

import (
	"context"
	"fmt"
	"personal_bot/app/iterable"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/pkg/logger"
	"sync"
	"time"
)

type TaskService struct {
	taskManager *Manager
	repo        *repository.TaskRepository
	hub         *subscriptionhub.Hub
	tasks       map[int64]tasks.Task
	mu          sync.Mutex
	iter        *iterable.Iterable
}

func NewTaskService(repo *repository.TaskRepository, hub *subscriptionhub.Hub, tManager *Manager) *TaskService {
	return &TaskService{
		taskManager: tManager,
		hub:         hub,
		repo:        repo,
		tasks:       map[int64]tasks.Task{},
		iter:        iterable.NewIterable(),
	}
}

func (ts *TaskService) Create(task tasks.Task) (tasks.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task.SetId(ts.iter.ID())
	task.SetTime(time.Now().Unix())

	ts.tasks[task.Id()] = task
	return task, nil
}

func (ts *TaskService) GetTaskWith(id int64) (tasks.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, ok := ts.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found with the id: %d", id)
	}
	return task, nil
}

func (ts *TaskService) GetTaskWithStrategyId(id int64) (tasks.Task, error) {
	for _, v := range ts.tasks {
		if v.Type() == "Sell" {
			continue
		}

		bt, ok := v.(*tasks.BuyTask)
		if !ok {
			return nil, fmt.Errorf("error whilst casting task: %d to buy task", id)
		}

		if *bt.StrategyId == id {
			return bt, nil
		}

	}

	return nil, fmt.Errorf("task not found with id: %d", id)
}

func (ts *TaskService) GetAllTasks() []tasks.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	allTasks := make([]tasks.Task, 0, len(ts.tasks))
	for _, val := range ts.tasks {
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

func (ts *TaskService) DeleteTask(id int64) (err error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	_, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	delete(ts.tasks, id)
	return nil
}

func (ts *TaskService) StartTask(id int64) (err error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task, ok := ts.tasks[id]

	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	if !task.State().TaskState.IsAbleToRun() {
		return fmt.Errorf("task is in terminal state - not able to run - State: %s", task.State())
	}

	task.SetState(tasks.State{TaskState: tasks.TaskValidating})
	//set to the next state after create
	if err := ts.taskManager.RunTask(task); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) StopTask(id int64) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	if err := ts.taskManager.StopTask(task); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) Subscribe(task tasks.Task) (*subscriptionhub.Subscription, error) {
	c, err := ts.hub.Subscribe(task)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (ts *TaskService) Unsubscribe(task tasks.Task) error {
	err := ts.hub.Unsubcribe(task)
	if err != nil {
		return err
	}
	return nil
}

// this will be called once on application load
func (ts *TaskService) LoadFromDB(ctx context.Context) error {
	tasksFromDb, err := ts.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	maxId := ts.repo.GetMaxId(ctx)
	ts.iter.SetIterable(maxId)

	logger.Information(tasksFromDb)

	for _, tdb := range tasksFromDb {
		ts.tasks[tdb.Id()] = tdb
	}

	return nil
}

// This will send all the tasks to the database before shutdown
func (ts *TaskService) Shutdown(ctx context.Context) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	deleteSuccess, err := ts.repo.DeleteAll(ctx)
	if err != nil {
		return err
	}

	if !deleteSuccess {
		return fmt.Errorf("failed to wipe table for insertion")
	}

	tasksToSave := make([]tasks.Task, 0, len(ts.tasks))
	for _, v := range ts.tasks {
		if ts.shouldPresistTask(v) {
			tasksToSave = append(tasksToSave, v)
		}
	}

	success, err := ts.repo.AddAllTasks(tasksToSave, ctx)
	if err != nil {
		return fmt.Errorf("error whilst shutting down: %w", err)
	}

	if !success {
		return fmt.Errorf("unsuccesful graceful shutdown in tasks service")
	}

	return nil
}

func (ts *TaskService) shouldPresistTask(t tasks.Task) bool {
	if sellT, ok := t.(*tasks.SellTask); ok {
		if sellT.Position_id != nil {
			return false
		}
	}

	return true
}
