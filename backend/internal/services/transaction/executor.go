package transaction

import (
	"context"
	"fmt"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction/buy"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction/sell"
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
	//TODO: we need to change this - pumpfun package will return the right transaction - so instead of switching task.type we'd switch on program
	switch t := task.(type) {
	case *tasks.BuyTask:
		return &buy.Transaction{BuyTask: t, PositionService: e.positionService}, nil
	case *tasks.SellTask:
		return &sell.Transaction{Task: t, PositionService: e.positionService}, nil
	}
	return nil, fmt.Errorf("no transaction found for task: %s", task.Type())
}

func (e *Executor) Execute(ctx context.Context, done chan struct{}, tx transaction.Transaction) {
	defer close(done)

	t := tx.GetTask()

	for {
		state, ok := e.steps[t.State().TaskState]

		// terminal state
		if !ok {
			return
		}

		// if there is ctx error/ cancellation then we set
		if err := ctx.Err(); err != nil {
			e.setStateAndPublish(tasks.TaskCancel, t)
			return
		}

		if err := state.Fn(ctx, tx); err != nil {
			if ctx.Err() != nil {
				e.setStateAndPublish(tasks.TaskCancel, t)
			} else {
				e.setStateAndPublish(state.OnError, t, err.Error())
			}
			return
		}

		e.setStateAndPublish(state.To, t)
	}
}
func (e *Executor) setStateAndPublish(newState tasks.TaskState, t tasks.Task, errMsg ...string) {
	state := tasks.State{TaskState: newState}
	if len(errMsg) > 0 && errMsg[0] != "" {
		state.Error = errMsg[0]
		t.SetMessage(errMsg[0])
	}
	t.SetState(state)
	e.publisher.PublishStateChange(t)
}
