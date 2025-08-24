package dto

import "pump_fun/internal/core/tasks"

type Subscribe struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type SubResponse struct {
	TaskEvent *tasks.TaskEvent `json:"task_event,omitempty"`
	Error     string           `json:"error,omitempty"`
}
