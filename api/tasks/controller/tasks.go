package controller

import (
	"pump_fun/api/mapper"
	"pump_fun/api/models"
	"pump_fun/internal/models/tasks"
	taskservice "pump_fun/internal/task_service"
)

type TaskController struct {
	TaskService *taskservice.TaskService
}

func (tc *TaskController) CreateTask(requestTask models.RequestTask) (tasks.Task, error) {
	// get the req struct
	// map to buy task or sell task depending
	newTask, err := mapper.MapRequestTaskToTask(&requestTask)
	if err != nil {
		return nil, err
	}
	// send to task service to create task
	createdTask, err := tc.TaskService.Create(newTask)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (tc *TaskController) GetTask(id string) (tasks.Task, error) {
	task, err := tc.TaskService.GetTaskWith(id)

	if err != nil {
		return nil, err
	}

	return *task, err
}

func (tc *TaskController) GetAllTasks() []tasks.Task {
	tasks := tc.TaskService.GetAllTasks()
	return tasks
}

func (tc *TaskController) UpdateTask(id string, reqTask models.RequestTask) (tasks.Task, error) {
	// we call update with => id + newTask
	updated, err := tc.TaskService.UpdateTask(id, reqTask)
	if err != nil {
		return nil, err
	}

	return *updated, nil
}

func (tc *TaskController) DeleteTask(id string) (err error) {
	err = tc.TaskService.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskController) TestEP() string {
	return "Test successful"
}
