package transaction

import (
	"fmt"
	"pump_fun/internal/handlers"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/transaction/buy"
	"pump_fun/internal/transaction/sell"
	"pump_fun/internal/transition"
	"pump_fun/pkg/logger"
)

type TransactionExecutor struct{}

func (th *TransactionExecutor) GetImplementation(task tasks.Task) Transaction {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return &buy.BuyTransaction{BuyTask: t}
	case *tasks.SellTask:
		return &sell.SellTransaction{Task: t}
	}

	return nil
}

func (th *TransactionExecutor) Execute(transaction Transaction) error {

	task := transaction.GetTask()
	ctx := task.Context()

	if task == nil {
		return fmt.Errorf("task wasn't set for the transaction")
	}

	err := handlers.ValidateStruct(task)
	if err != nil {
		transition.AutoTransitionTask(task, err)
		return err
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
			return err
		}

		err := step()
		transition.AutoTransitionTask(task, err)
		if err != nil {
			return err
		}
	}

	return nil

}
