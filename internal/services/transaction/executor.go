package transaction

import (
	"context"
	"fmt"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/core/validator"
	"pump_fun/internal/services/state/transition"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/internal/solana/programs/pumpfun/transaction/buy"
	"pump_fun/internal/solana/programs/pumpfun/transaction/sell"
)

type Executor struct {
	subhub *subscriptionhub.Hub
}

func (e *Executor) New(subhub *subscriptionhub.Hub) {
	e.subhub = subhub
}

func (e *Executor) GetImplementation(task tasks.Task) (Transaction, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return &buy.Transaction{BuyTask: t}, nil
	case *tasks.SellTask:
		return &sell.Transaction{Task: t}, nil
	}

	return nil, fmt.Errorf("no transaction found for the task: %s", task.Type())
}

func (e *Executor) Execute(done chan struct{}, transaction Transaction, ctx context.Context) {

	task := transaction.GetTask()

	reporter := subscriptionhub.TaskReporter{}
	reporter.New(task, e.subhub)

	err := validator.ValidateStruct(task)
	if err != nil {
		e.transitionAndPublishTask(task, err)
		close(done)
		return
	}

	err = transition.AutoTransitionTask(task, nil) //from running to next step
	if err != nil {
		e.transitionAndPublishTask(task, err)
		close(done)
		return
	}

	steps := []func(ctx context.Context, reporter subscriptionhub.TaskReporter) error{
		transaction.BuildInstructions,
		transaction.BuildTransaction,
		transaction.SendTransaction,
		transaction.ConfirmTransaction,
	}

	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			e.transitionAndPublishTask(task, err)
			close(done)
			return
		}

		err := step(ctx, reporter)
		e.transitionAndPublishTask(task, err)
		if err != nil {
			close(done)
			return
		}
	}

	close(done)
}

func (e *Executor) transitionAndPublishTask(t tasks.Task, err error) {
	err = transition.AutoTransitionTask(t, err)
	if err != nil {
		return
	}

	e.subhub.PublishStateChange(t)
}
