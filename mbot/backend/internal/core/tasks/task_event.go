package tasks

type TaskEvent struct {
	TaskId     int64     `json:"task_id"`
	StrategyId int64     `json:"strategy_id"`
	State      StateDto  `json:"state"`
	Time       string    `json:"time"`
	Message    string    `json:"message"`
	EventType  EventType `json:"event_type"`
}

type StateDto struct {
	TaskState string `json:"task_state"`
	Error     string `json:"error"`
}

type EventType string

const (
	StateUpdate     = "StatusUpdate"
	ProgressMessage = "ProgressMessage"
)
