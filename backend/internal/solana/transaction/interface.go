package transaction

import (
	"context"
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
)

type Transaction interface {
	BuildInstructions(ctx context.Context, publisher subscriptionhub.Publisher) error
	BuildTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error
	SendTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error
	ConfirmTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error
	UpdatePosition(ctx context.Context, publisher subscriptionhub.Publisher) (tokenAmount, solAmount float64, pos *position.Position, err error)
	GetTask() tasks.Task
	GetSignature() string
	GetAddressForURL() (string, error)
}
