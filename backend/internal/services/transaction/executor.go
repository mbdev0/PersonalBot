package transaction

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/core/validator"
	"personal_bot/internal/services/position"
	"personal_bot/internal/services/state/transition"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	"personal_bot/internal/solana/programs/pumpfun/transaction/buy"
	"personal_bot/internal/solana/programs/pumpfun/transaction/sell"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

type Executor struct {
	subhub          *subscriptionhub.Hub
	positionService *position.Service
}

func (e *Executor) New(subhub *subscriptionhub.Hub, posService *position.Service) {
	e.subhub = subhub
	e.positionService = posService

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
	defer close(done)

	task := transaction.GetTask()
	reporter := subscriptionhub.TaskReporter{}
	reporter.New(task, e.subhub)

	err := validator.ValidateStruct(task)
	if err != nil {
		e.transitionAndPublishTask(task, err)
		return
	}

	e.handlePositionOnSell(task, ctx)

	err = transition.AutoTransitionTask(task, nil) //from running to next step
	if err != nil {
		e.transitionAndPublishTask(task, err)
		return
	}

	steps := []func(ctx context.Context, reporter subscriptionhub.TaskReporter) error{

		func(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
			return transaction.BuildInstructionsWithPosition(ctx, reporter, e.positionService)
		},
		transaction.BuildTransaction,
		transaction.SendTransaction,
		transaction.ConfirmTransaction,
	}

	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			e.transitionAndPublishTask(task, err)
			return
		}

		err := step(ctx, reporter)
		e.transitionAndPublishTask(task, err)
		if err != nil {
			return
		}
	}
	if task.State().TaskState == tasks.TaskDone {
		err := e.updatePositionOnCompleted(task, transaction, ctx)
		if err != nil {
			task.SetState(tasks.State{TaskState: tasks.TaskDone, Error: fmt.Sprintf("unable to create position: %s", err.Error())})
		}
	}
}

func (e *Executor) transitionAndPublishTask(t tasks.Task, err error) {
	err = transition.AutoTransitionTask(t, err)
	if err != nil {
		return
	}

	e.subhub.PublishStateChange(t)
}

func (e *Executor) handlePositionOnSell(task tasks.Task, ctx context.Context) {
	st, ok := task.(*tasks.SellTask)
	if !ok || st.Position_id != nil {
		return
	}

	_, exists := e.positionService.FindPositionIfExists(st.Token, st.Wallet.PublicKey())
	if exists {
		return
	}

	var err error
	var ata solana.PublicKey

	isNewToken, err := instructions.IsTokenAccountNew(st.Token, ctx)
	if err != nil {
		e.transitionAndPublishTask(task, err)
	}

	if isNewToken {
		ata, _, err = pda.FindToken2022AssociatedTokenAddress(st.Wallet.PublicKey(), st.Token)
	} else {
		ata, _, err = pda.FindTokenAssociatedTokenAddress(st.Wallet.PublicKey(), st.Token)
	}

	if err != nil {
		e.transitionAndPublishTask(task, err)
		return
	}

	tokens, err := client.GetTokenAccountBalance(ata, ctx)
	if err != nil {
		e.transitionAndPublishTask(task, err)
		return
	}

	e.positionService.ReportBuy(ctx, st.Id(), st.Token, st.Wallet.PublicKey(), new(big.Float).SetUint64(*tokens), new(big.Float).SetFloat64(0))

}
func (e *Executor) updatePositionOnCompleted(task tasks.Task, transaction Transaction, ctx context.Context) error {
	tokenAmount, solAmount, err := transaction.ExtractTokenAndSolFromTx(transaction.GetSignature(), ctx)
	if err != nil {
		logger.Error("error: ", err)
		e.transitionAndPublishTask(task, err)
		return nil
	}

	switch t := task.(type) {
	case *tasks.BuyTask:
		e.positionService.ReportBuy(ctx, t.Id(), t.Token, t.Wallet.PublicKey(), new(big.Float).SetFloat64(tokenAmount), new(big.Float).SetFloat64(solAmount))
	case *tasks.SellTask:
		err := e.handleSellTaskReporting(t, tokenAmount, solAmount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) handleSellTaskReporting(t *tasks.SellTask, tokenAmount float64, solAmount float64) error {
	tokensSold := new(big.Float).SetFloat64(tokenAmount)
	solReceived := new(big.Float).SetFloat64(solAmount)

	if t.Position_id == nil {
		pos, _ := e.positionService.FindPositionIfExists(t.Token, t.Wallet.PublicKey())
		err := e.positionService.ReportSell(pos.PositionId, tokensSold, solReceived)
		if err != nil {
			return err
		}
	} else {
		err := e.positionService.ReportSell(*t.Position_id, tokensSold, solReceived)
		if err != nil {
			return err
		}
	}
	return nil
}
