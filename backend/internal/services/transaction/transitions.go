package transaction

import (
	"context"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/transaction"
)

type State struct {
	From    tasks.TaskState
	To      tasks.TaskState
	Fn      func(ctx context.Context, t transaction.Transaction) error
	OnError tasks.TaskState
}

type Transitions map[tasks.TaskState]State
