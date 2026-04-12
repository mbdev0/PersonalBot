package transaction

import (
	"context"
	"fmt"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/programs/pumpfun/transaction/buy"
	"personal_bot/internal/solana/programs/pumpfun/transaction/sell"
	"personal_bot/internal/solana/transaction"
)

type Executor struct {
	publisher       subscriptionhub.Publisher
	positionService *position.Service
	steps           Transitions
}

func NewExecutor(publisher subscriptionhub.Publisher, posService *position.Service, steps Transitions) *Executor {
	return &Executor{
		publisher:       publisher,
		positionService: posService,
		steps:           steps,
	}
}

func (e *Executor) GetImplementation(task tasks.Task) (transaction.Transaction, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return &buy.Transaction{BuyTask: t, PositionService: e.positionService}, nil
	case *tasks.SellTask:
		return &sell.Transaction{Task: t, PositionService: e.positionService}, nil
	}

	return nil, fmt.Errorf("no transaction found for the task: %s", task.Type())
}

func (e *Executor) Execute(ctx context.Context, done chan struct{}, transaction transaction.Transaction) {
	defer close(done)

	t := transaction.GetTask()

	for {
		state, ok := e.steps[t.State().TaskState]

		//terminal state
		if !ok {
			return
		}

		// if there is ctx error/ cancellation then we set
		if err := ctx.Err(); err != nil {
			e.setStateAndPublish(tasks.TaskCancel, t)
			return
		}

		if err := state.Fn(ctx, transaction); err != nil {
			if ctx.Err() != nil {
				e.setStateAndPublish(tasks.TaskCancel, t)
			} else {
				e.setStateAndPublish(state.OnError, t)
			}
			return
		}

		e.setStateAndPublish(state.To, t)

	}
}

func (e *Executor) setStateAndPublish(newState tasks.TaskState, t tasks.Task) {
	t.SetState(tasks.State{TaskState: newState})
	e.publisher.PublishStateChange(t)
}
