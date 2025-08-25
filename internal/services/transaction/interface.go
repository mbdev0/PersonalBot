package transaction

import (
	"pump_fun/internal/core/tasks"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
)

type Transaction interface {
	BuildInstructions(reporter subscriptionhub.TaskReporter) error
	BuildTransaction(reporter subscriptionhub.TaskReporter) error
	SendTransaction(reporter subscriptionhub.TaskReporter) error
	ConfirmTransaction(reporter subscriptionhub.TaskReporter) error
	GetTask() tasks.Task
}
