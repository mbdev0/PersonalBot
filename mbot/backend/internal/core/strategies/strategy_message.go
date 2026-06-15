package strategies

import "personal_bot/backend/internal/core/tasks"

type StrategyMessage struct {
	Id      int64      `json:"id"`
	Event   EventType  `json:"event"`
	Task    tasks.Task `json:"task,omitempty"`
	State   *string    `json:"state,omitempty"`
	Message *string    `json:"message,omitempty"`
}

type EventType string

const (
	TaskCreation  EventType = "TASK_CREATION"
	StatusUpdate  EventType = "STATUS_UPDATE"
	MessageUpdate EventType = "MESSAGE_UPDATE"
)
