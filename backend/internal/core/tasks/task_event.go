package tasks

type TaskEvent struct {
	TaskId    string
	State     State
	Time      string
	Message   string
	EventType EventType
}

type EventType string

const (
	StateUpdate     = "StatusUpdate"
	ProgressMessage = "ProgressMessage"
)
