package dto

import (
	"pump_fun/internal/core/position"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/core/tasks"
)

type TaskSubscribe struct {
	Type SubType `json:"type"`
	Id   int64   `json:"id"`
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
	StrategyMessage *strategies.StrategyMessage `json:"strategy_msg,omitempty"`
	Error           string                      `json:"error,omitempty"`
}

type SubType string

const (
	Subscribe   SubType = "Subscribe"
	Unsubscribe SubType = "Unsubscribe"
)
