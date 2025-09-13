package transaction

import (
	"context"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/services/position"
	subscriptionhub "pump_fun/internal/services/subscription_hub"

	"github.com/gagliardetto/solana-go"
)

type Transaction interface {
	// BuildInstructions(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	BuildInstructionsWithPosition(ctx context.Context, reporter subscriptionhub.TaskReporter, ps *position.Service) error
	BuildTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	SendTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	ConfirmTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	ExtractTokenAndSolFromTx(signature solana.Signature, ctx context.Context) (tokenAmount float64, solAmount float64, err error)
	GetTask() tasks.Task
	GetSignature() solana.Signature
}
