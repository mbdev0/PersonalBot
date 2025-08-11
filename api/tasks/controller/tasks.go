package controller

import (
	"pump_fun/api/mapper"
	"pump_fun/api/models"
	"pump_fun/internal/core/tasks"
	taskservice "pump_fun/internal/services/task_service"
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

func (tc *TaskController) GetTask(id string) (*models.ResponseTask, error) {
	task, err := tc.TaskService.GetTaskWith(id)
	if err != nil {
		return nil, err
	}

	response, err := mapper.MapTaskToReponseTask(task)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (tc *TaskController) GetAllTasks() ([]models.ResponseTask, error) {
	tasks := tc.TaskService.GetAllTasks()

	response := make([]models.ResponseTask, 0, len(tasks))
	for _, task := range tasks {
		responseObj, err := mapper.MapTaskToReponseTask(task)
		if err != nil {
			return nil, err
		}
		response = append(response, *responseObj)
	}

	return response, nil
}

func (tc *TaskController) UpdateTask(id string, reqTask models.RequestTask) (tasks.Task, error) {
	// we call update with => id + newTask
	updated, err := tc.TaskService.UpdateTask(id, reqTask)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (tc *TaskController) DeleteTask(id string) (err error) {
	err = tc.TaskService.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskController) TransitionTask(id string, newState string) (err error) {
	state, err := tasks.ParseStateString(newState)
	if err != nil {
		return err
	}

	err = tc.TaskService.TransitionTask(id, state)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskController) TestEP() string {
	return "Test successful"
}
