package cryptostates

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/core/validator"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/transaction"
)

type Transitions map[tasks.TaskState]State

type Dependencies struct {
	Publisher       subscriptionhub.Publisher
	PositionService *position.Service
}

func Build(deps *Dependencies) Transitions {
	return Transitions{
		tasks.TaskValidating: {
			From: tasks.TaskValidating,
			To:   tasks.TxInstructionBuild,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return validator.ValidateStruct(t.GetTask())
			},
			OnError: tasks.TaskValidationFail,
		},
		tasks.TxInstructionBuild: {
			From: tasks.TxInstructionBuild,
			To:   tasks.TxBuild,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return t.BuildInstructionsWithPosition(ctx, deps.Publisher, deps.PositionService)
			},
			OnError: tasks.TxInstructionBuildFail,
		},
		tasks.TxBuild: {
			From: tasks.TxBuild,
			To:   tasks.TxSend,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return t.BuildTransaction(ctx, deps.Publisher)
			},
			OnError: tasks.TxBuildFail,
		},
		tasks.TxSend: {
			From: tasks.TxSend,
			To:   tasks.TxConfirm,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return t.SendTransaction(ctx, deps.Publisher)
			},
			OnError: tasks.TxSendFail,
		},
		tasks.TxConfirm: {
			From: tasks.TxConfirm,
			To:   tasks.TaskUpdatingPosition,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return t.ConfirmTransaction(ctx, deps.Publisher)
			},
			OnError: tasks.TxConfirmFail,
		},
		tasks.TaskUpdatingPosition: {
			From: tasks.TaskUpdatingPosition,
			To:   tasks.TaskDone,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return updatePositionOnCompleted(t.GetTask(), t, ctx, deps)
			},
			OnError: tasks.TaskUpdatingPositionFail,
		},
	}
}

func updatePositionOnCompleted(task tasks.Task, transaction transaction.Transaction, ctx context.Context, deps *Dependencies) error {
	tokenAmount, solAmount, err := transaction.ExtractTokenAndSolFromTx(transaction.GetSignature(), ctx)
	if err != nil {
		return err
	}

	switch t := task.(type) {
	case *tasks.BuyTask:
		deps.PositionService.ReportBuy(ctx, t.Id(), t.Token, t.Wallet.PublicKey(), new(big.Float).SetFloat64(tokenAmount), new(big.Float).SetFloat64(solAmount))
	case *tasks.SellTask:
		return handleSellTaskReporting(t, tokenAmount, solAmount, deps)
	}

	return nil
}

func handleSellTaskReporting(t *tasks.SellTask, tokenAmount float64, solAmount float64, deps *Dependencies) error {
	tokensSold := new(big.Float).SetFloat64(tokenAmount)
	solReceived := new(big.Float).SetFloat64(solAmount)

	if t.Position_id == nil {
		pos, exists := deps.PositionService.FindPositionIfExists(t.Token, t.Wallet.PublicKey())
		if !exists {
			return fmt.Errorf("position not found for sell task %d: ", t.Id())
		}

		err := deps.PositionService.ReportSell(pos.PositionId, tokensSold, solReceived)
		if err != nil {
			return err
		}
	} else {
		err := deps.PositionService.ReportSell(*t.Position_id, tokensSold, solReceived)
		if err != nil {
			return err
		}
	}
	return nil
}
