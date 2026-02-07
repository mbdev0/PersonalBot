package transaction

import (
	"context"
	"fmt"
	"personal_bot/internal/core/tasks"
	cryptostates "personal_bot/internal/core/tasks/crypto_states"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/programs/pumpfun/transaction/buy"
	"personal_bot/internal/solana/programs/pumpfun/transaction/sell"
	"personal_bot/internal/solana/transaction"
)

type Executor struct {
	publisher       subscriptionhub.Publisher
	positionService *position.Service
	steps           cryptostates.Transitions
}

func NewExecutor(publisher subscriptionhub.Publisher, posService *position.Service, steps cryptostates.Transitions) *Executor {
	return &Executor{
		publisher:       publisher,
		positionService: posService,
		steps:           steps,
	}
}

func (e *Executor) GetImplementation(task tasks.Task) (transaction.Transaction, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return &buy.Transaction{BuyTask: t}, nil
	case *tasks.SellTask:
		return &sell.Transaction{Task: t}, nil
	}

	return nil, fmt.Errorf("no transaction found for the task: %s", task.Type())
}

func (e *Executor) Execute(done chan struct{}, transaction transaction.Transaction, ctx context.Context) {
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
			continue
		}

		if err := state.Fn(ctx, transaction); err != nil {
			//transition to on error
			e.setStateAndPublish(state.OnError, t)
		}

		e.setStateAndPublish(state.To, t)

	}
}

func (e *Executor) setStateAndPublish(newState tasks.TaskState, t tasks.Task) {
	t.SetState(tasks.State{TaskState: newState})
	e.publisher.PublishStateChange(t)
}
