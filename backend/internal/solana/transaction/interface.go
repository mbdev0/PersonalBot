package transaction

import (
	"context"
	"personal_bot/internal/core/position"
)

type Transaction interface {
	BuildInstructions(ctx context.Context) error
	BuildTransaction(ctx context.Context) error
	SendTransaction(ctx context.Context) error
	ConfirmTransaction(ctx context.Context) error
	UpdatePosition(ctx context.Context) (tokenAmount, solAmount float64, pos *position.Position, err error)
	GetSignature() string
	GetAddressForURL() (string, error)
}
