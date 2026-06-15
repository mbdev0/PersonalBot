package subscriptionhub

import "personal_bot/backend/internal/core/tasks"

type Publisher interface {
	PublishMessage(task tasks.Task, message string)
	PublishStateChange(task tasks.Task)
}
