package tasks

import "pump_fun/internal/core/position"

type TaskEvent struct {
	TaskId          string
	State           State
	Time            string
	Message         string
	EventType       EventType
	PositionDetails *position.PositionMessage
}

type EventType string

const (
	StateUpdate     = "StatusUpdate"
	ProgressMessage = "ProgressMessage"
	PositionUpdate  = "PositionUpdate"
)
