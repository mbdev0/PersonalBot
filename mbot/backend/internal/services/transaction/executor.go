package transaction

import (
	"context"
	"fmt"
	"personal_bot/backend/internal/core/program"
	"personal_bot/backend/internal/core/tasks"
	"personal_bot/backend/internal/services/position"
	subscriptionhub "personal_bot/backend/internal/services/subscription_hub"
	"personal_bot/backend/internal/solana/transaction"
	"time"
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

func (e *Executor) GetImplementation(ctx context.Context, task tasks.Task) (transaction.Transaction, error) {
	program := program.Resolve(task.Program())
	switch t := task.(type) {
	case *tasks.BuyTask:
		return program.NewBuyTransaction(t, e.positionService, e.publisher), nil
	case *tasks.SellTask:
		return program.NewSellTransaction(t, e.positionService, e.publisher), nil
	default:
		return nil, fmt.Errorf("unable to find program: %s", task.Program())
	}
}

func (e *Executor) Execute(ctx context.Context, done chan struct{}, tx transaction.Transaction, task tasks.Task) {
	defer close(done)

	retries := uint16(0)
	t, isRetryable := task.(tasks.RetryableTask)
	if isRetryable {
		retries = t.Retries()
	}

	for {
		state, ok := e.steps[task.State().TaskState]

		// terminal state
		if !ok {
			return
		}

		// if there is ctx error/ cancellation then we set
		if err := ctx.Err(); err != nil {
			e.setStateAndPublish(tasks.TaskCancel, task)
			return
		}

		if err := state.Fn(ctx, task, tx); err != nil {
			if ctx.Err() != nil {
				e.setStateAndPublish(tasks.TaskCancel, task)
			} else {
				if isRetryable && retries > 0 {
					retries--
					time.Sleep(time.Duration(t.RetriesDelayMS()) * time.Millisecond)
					e.setStateAndPublish(t.RetryFrom(), task)
					continue
				}
				e.setStateAndPublish(state.OnError, task, err.Error())
			}
			return
		}

		e.setStateAndPublish(state.To, task)
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
