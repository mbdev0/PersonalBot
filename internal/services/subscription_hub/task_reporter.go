package subscriptionhub

import (
	"pump_fun/internal/core/tasks"
)

type TaskReporter struct {
	publisher Publisher
	task      tasks.Task
}

func (tr *TaskReporter) New(t tasks.Task, p Publisher) {
	tr.publisher = p
	tr.task = t
}

func (tr *TaskReporter) Report(message string) {
	tr.publisher.PublishMessage(tr.task, message)
}
