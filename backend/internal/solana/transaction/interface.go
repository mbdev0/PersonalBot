package transaction

import (
	"context"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"

	"github.com/gagliardetto/solana-go"
)

type Transaction interface {
	// BuildInstructions(ctx context.Context, reporter subscriptionhub.TaskReporter) error
	BuildInstructionsWithPosition(ctx context.Context, publisher subscriptionhub.Publisher, ps *position.Service) error
	BuildTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error
	SendTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error
	ConfirmTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error
	ExtractTokenAndSolFromTx(ctx context.Context, signature solana.Signature) (tokenAmount float64, solAmount float64, err error)
	GetTask() tasks.Task
	GetSignature() solana.Signature
}
