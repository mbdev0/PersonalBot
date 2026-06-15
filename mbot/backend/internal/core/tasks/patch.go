package tasks

type TaskPatch interface {
	ApplyTo(t Task) error
}
