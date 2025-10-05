package subscriptionhub

import "pump_fun/internal/core/tasks"

type Publisher interface {
	PublishMessage(task tasks.Task, message string)
	PublishStateChange(task tasks.Task)
}
