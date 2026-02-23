package tasks

type State struct {
	TaskState TaskState `json:"task_state"`
	Error     string    `json:"error"`
}
