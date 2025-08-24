package tasks

type TaskEvent struct {
	TaskId string
	State  State
	Time   string
}
