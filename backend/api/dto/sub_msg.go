package dto

import (
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
)

type TaskSubscribe struct {
	Type       SubType `json:"type"`
	Id         *int64  `json:"id,omitempty"`
	StrategyId *int64  `json:"strategy_id,omitempty"`
}

type TaskSubResponse struct {
	TaskEvent *tasks.TaskEvent `json:"task_event,omitempty"`
	Error     string           `json:"error,omitempty"`
}

type PositionSubscribe struct {
	Type SubType `json:"type"`
	Id   int64   `json:"id"`
}

type PositionResponse struct {
	PositionMessage *position.PositionMessage `json:"position_msg,omitempty"`
	Error           string                    `json:"error,omitempty"`
}

type StrategySubscribe struct {
	Type SubType `json:"type"`
	Id   int64   `json:"id"`
}

type StrategyResponse struct {
	StrategyMessage StrategyMessageResponse `json:"strategy_msg,omitempty"`
	Error           string                  `json:"error,omitempty"`
}

type StrategyMessageResponse struct {
	Id    int64        `json:"id"`
	Event string       `json:"event"`
	Task  ResponseTask `json:"task"`
}

type SubType string

const (
	Subscribe   SubType = "Subscribe"
	Unsubscribe SubType = "Unsubscribe"
)
