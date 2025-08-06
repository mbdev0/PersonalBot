package models

import "context"

type CancelToken struct {
	CancellationContext context.Context
	CancellationFunc    context.CancelFunc
}
