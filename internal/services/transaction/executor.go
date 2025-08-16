package transaction

import (
	"fmt"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/core/validator"
	"pump_fun/internal/services/state/transition"
	"pump_fun/internal/solana/programs/pumpfun/transaction/buy"
	"pump_fun/internal/solana/programs/pumpfun/transaction/sell"
	"pump_fun/pkg/logger"
)

type Executor struct{}

func (th *Executor) GetImplementation(task tasks.Task) (Transaction, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return &buy.BuyTransaction{BuyTask: t}, nil
	case *tasks.SellTask:
		return &sell.SellTransaction{Task: t}, nil
	}

	return nil, fmt.Errorf("no transaction found for the task: %s", task.Type())
}

func (th *Executor) Execute(done chan bool, transaction Transaction) {

	task := transaction.GetTask()
	ctx := task.Ctx()

	if task == nil {
		transition.AutoTransitionTask(task, fmt.Errorf("no task set in transaction"))
		done <- true
		return
	}

	err := validator.ValidateStruct(task)
	if err != nil {
		transition.AutoTransitionTask(task, err)
		done <- true
		return
	}

	transition.AutoTransitionTask(task, nil) //from running to next step

	steps := []func() error{
		transaction.BuildInstructions,
		transaction.BuildTransaction,
		transaction.SendTransaction,
		transaction.ConfirmTransaction,
	}

	logger.Information(steps)

	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			transition.AutoTransitionTask(task, err)
			done <- true
			return
		}

		err := step()
		transition.AutoTransitionTask(task, err)
		if err != nil {
			done <- true
			return
		}
	}

	done <- true
}
