package tasks

type TaskEvent struct {
	TaskId     int64
	StrategyId int64
	State      State
	Time       string
	Message    string
	EventType  EventType
}

type EventType string

const (
	StateUpdate     = "StatusUpdate"
	ProgressMessage = "ProgressMessage"
)
