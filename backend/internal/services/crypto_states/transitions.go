package cryptostates

import (
	"context"
	"personal_bot/internal/core/constants"
	notifierModel "personal_bot/internal/core/notifier"
	positionModel "personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/core/validator"
	"personal_bot/internal/services/notifier"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	taskservice "personal_bot/internal/services/task_service"
	"personal_bot/internal/services/trading"
	transactionModel "personal_bot/internal/services/transaction"
	"personal_bot/internal/solana/transaction"
	"personal_bot/pkg/logger"
)

type Dependencies struct {
	Publisher       subscriptionhub.Publisher
	PositionService *position.Service
	Notifier        *notifier.DiscordNotifier
	TradingService  *trading.Service
	TaskService     *taskservice.TaskService
}

func Build(deps *Dependencies) transactionModel.Transitions {
	return transactionModel.Transitions{
		tasks.TaskValidating: {
			To:      tasks.TxInstructionBuild,
			OnError: tasks.TaskValidationFail,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return withNotify(ctx, deps, func() error {
					return validator.ValidateStruct(t.GetTask())
				}, t.GetTask(), t)
			},
		},
		tasks.TxInstructionBuild: {
			To:      tasks.TxBuild,
			OnError: tasks.TxInstructionBuildFail,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return withNotify(ctx, deps, func() error {
					return t.BuildInstructions(ctx, deps.Publisher)
				}, t.GetTask(), t)
			},
		},
		tasks.TxBuild: {
			To:      tasks.TxSend,
			OnError: tasks.TxBuildFail,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return withNotify(ctx, deps, func() error {
					return t.BuildTransaction(ctx, deps.Publisher)
				}, t.GetTask(), t)
			},
		},
		tasks.TxSend: {
			To:      tasks.TxConfirm,
			OnError: tasks.TxSendFail,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return withNotify(ctx, deps, func() error {
					return t.SendTransaction(ctx, deps.Publisher)
				}, t.GetTask(), t)
			},
		},
		tasks.TxConfirm: {
			To:      tasks.TaskUpdatingPosition,
			OnError: tasks.TxConfirmFail,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return withNotify(ctx, deps, func() error {
					return t.ConfirmTransaction(ctx, deps.Publisher)
				}, t.GetTask(), t)
			},
		},
		tasks.TaskUpdatingPosition: {
			To:      tasks.TaskDone,
			OnError: tasks.TaskUpdatingPositionFail,
			Fn: func(ctx context.Context, t transaction.Transaction) error {
				return withNotify(ctx, deps, func() error {
					tokenAmount, solAmount, pos, err := t.UpdatePosition(ctx, deps.Publisher)
					if err != nil {
						return err
					}
					notifyCompletion(t.GetTask(), t, tokenAmount, solAmount, pos, deps)
					return nil
				}, t.GetTask(), t)
			},
		},
	}
}

func withNotify(ctx context.Context, deps *Dependencies, fn func() error, task tasks.Task, t transaction.Transaction) error {
	err := fn()
	if err != nil && ctx.Err() == nil {
		if notifyErr := notifyError(err.Error(), deps, task, t); notifyErr != nil {
			logger.Error(notifyErr)
		}
	}
	return err
}

func notifyCompletion(task tasks.Task, transaction transaction.Transaction, tokenAmount, solAmount float64, pos *positionModel.Position, deps *Dependencies) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		notifyBuy(t, transaction, tokenAmount, solAmount, deps)
	case *tasks.SellTask:
		notifySell(t, transaction, pos, tokenAmount, solAmount, deps)
	}
}

func notifyBuy(t *tasks.BuyTask, transaction transaction.Transaction, tokenAmount float64, solAmount float64, deps *Dependencies) {
	addressForAxiom, err := transaction.GetAddressForURL()
	if err != nil {
		logger.Error(err)
	}

	tradingTask, err := deps.TradingService.GetBy(*t.StrategyId)
	if err != nil {
		logger.Error(err)
		return
	}

	err = deps.Notifier.SendSuccessBuy(notifierModel.BuyNotifierPayload{
		TaskType:        string(tradingTask.StrategyType()),
		TaskId:          t.Id(),
		StrategyId:      t.StrategyId,
		TokensBought:    tokenAmount / constants.TokenAmountDecimals,
		AmountPaid:      solAmount / constants.LamportsConversion,
		TxHash:          transaction.GetSignature(),
		WalletAddress:   t.Wallet.PublicKey().String(),
		TokenAddress:    t.Token.String(),
		AddressForAxiom: addressForAxiom,
	})
	if err != nil {
		logger.Error(err)
	}
}

func notifySell(t *tasks.SellTask, transaction transaction.Transaction, pos *positionModel.Position, tokensSold float64, solReceived float64, deps *Dependencies) {
	addressForAxiom, err := transaction.GetAddressForURL()
	if err != nil {
		logger.Error(err)
	}

	if pos.TokenRemaining == nil {
		logger.Information("tokens remaning is nil")
		logger.Information(*pos)
	}

	tokensRemainingRaw, _ := pos.TokenRemaining.Float64()
	currentProfitRaw, _ := pos.FinalizedProfit.Float64()

	payload := notifierModel.SellNotifierPayload{
		TaskId:          t.Id(),
		StrategyId:      t.StrategyId,
		TxHash:          transaction.GetSignature(),
		TokensSold:      tokensSold / constants.TokenAmountDecimals,
		TokensRemaining: tokensRemainingRaw / constants.TokenAmountDecimals,
		CurrentProfit:   currentProfitRaw / constants.LamportsConversion,
		AmountSold:      solReceived / constants.LamportsConversion,
		WalletAddress:   t.Wallet.PublicKey().String(),
		TokenAddress:    t.Token.String(),
		AddressForAxiom: addressForAxiom,
	}

	if t.StrategyId != nil {
		strategyTask, err := deps.TradingService.GetBy(*t.StrategyId)
		if err != nil {
			logger.Error(err)
			return
		}
		payload.TaskType = string(strategyTask.StrategyType())
	} else {
		payload.TaskType = "SellTask"
	}

	err = deps.Notifier.SendSuccessSell(payload)

	if err != nil {
		logger.Error(err)
	}
}

func notifyError(errorMessage string, deps *Dependencies, task tasks.Task, transaction transaction.Transaction) error {
	if task.State().TaskState.IsCancel() {
		return nil
	}

	var task_type string

	if task.GetStrategyId() != nil {
		strategy, err := deps.TradingService.GetBy(*task.GetStrategyId())
		if err != nil {
			return err
		}
		task_type = string(strategy.StrategyType())
	} else {
		task, err := deps.TaskService.GetTaskWith(task.Id())
		if err != nil {
			return err
		}
		task_type = task.Type()
	}

	addressForAxiom, err := transaction.GetAddressForURL()
	if err != nil {
		return err
	}

	err = deps.Notifier.SendFailure(notifierModel.ErrorNotifierPayload{
		TaskType:        task_type,
		TaskId:          task.Id(),
		StrategyId:      task.GetStrategyId(),
		Error:           errorMessage,
		TxHash:          transaction.GetSignature(),
		WalletAddress:   task.GetWallet(),
		TokenAddress:    task.GetToken(),
		AddressForAxiom: addressForAxiom,
	})

	return err
}
