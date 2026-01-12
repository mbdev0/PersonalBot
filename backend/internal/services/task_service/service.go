package taskservice

import (
	"context"
	"fmt"
	"maps"
	"pump_fun/infrastructure/persistence/repository"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/services/state"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/pkg/logger"
	"slices"
	"sync"
)

type TaskService struct {
	stateMachine *state.Machine
	stateManager *state.Manager
	repo         *repository.TaskRepository
	hub          *subscriptionhub.Hub
	tasks        map[int64]tasks.Task
	mu           sync.Mutex
}

func NewTaskService(sm *state.Machine, stateManager *state.Manager, repo *repository.TaskRepository, hub *subscriptionhub.Hub) *TaskService {
	return &TaskService{
		hub:          hub,
		stateMachine: sm,
		stateManager: stateManager,
		repo:         repo,
		tasks:        map[int64]tasks.Task{},
	}
}

func (ts *TaskService) Create(task tasks.Task) (tasks.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
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

func (ts *TaskService) TransitionTask(id int64, newState tasks.TaskState) (err error) {
	// in here we'll manage changing state
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("Task not found with the id: %d", id)
	}

	err = ts.stateMachine.Transition(task, newState)
	ts.hub.PublishStateChange(task)

	if err != nil {
		return fmt.Errorf("transition failed for task %d with error: %w", task.Id(), err)
	}

	err = ts.stateManager.ExecuteAction(newState, task)
	if err != nil {
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

	tasksSlice := slices.Collect(maps.Values(ts.tasks))

	success, err := ts.repo.AddAllTasks(tasksSlice, ctx)
	if err != nil {
		return fmt.Errorf("error whilst shutting down: %w", err)
	}

	if !success {
		return fmt.Errorf("unsuccesful graceful shutdown in tasks service")
	}

	return nil
}
