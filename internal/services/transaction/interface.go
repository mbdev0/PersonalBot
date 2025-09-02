package transaction

import (
	"context"
	"pump_fun/internal/core/tasks"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
)

type Transaction interface {
	BuildInstructions(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	BuildTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	SendTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	ConfirmTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	GetTask() tasks.Task
}
