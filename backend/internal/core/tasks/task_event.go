package tasks

type TaskEvent struct {
	TaskId     int64     `json:"task_id"`
	StrategyId int64     `json:"strategy_id"`
	State      State     `json:"state"`
	Time       string    `json:"time"`
	Message    string    `json:"message"`
	EventType  EventType `json:"event_type"`
}

type EventType string

const (
	StateUpdate     = "StatusUpdate"
	ProgressMessage = "ProgressMessage"
)
