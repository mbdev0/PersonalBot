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
	"personal_bot/internal/services/trading"
	transactionModel "personal_bot/internal/services/transaction"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	"personal_bot/internal/solana/transaction"
	"personal_bot/pkg/logger"
)

type Dependencies struct {
	Publisher       subscriptionhub.Publisher
	PositionService *position.Service
	Notifier        *notifier.DiscordNotifier
	TradingService  *trading.Service
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
					notifyCompletion(t.GetTask(), t, tokenAmount, solAmount, pos, deps)
					return err
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
	bondingCurve, err := pda.GetBondingCurveAddress(t.Token.String())
	if err != nil {
		logger.Error(err)
	}

	tradingTask, err := deps.TradingService.GetBy(*t.StrategyId)
	if err != nil {
		logger.Error(err)
		return
	}

	err = deps.Notifier.SendSuccessBuy(notifierModel.BuyNotifierPayload{
		TaskType:      string(tradingTask.StrategyType()),
		TaskId:        t.Id(),
		StrategyId:    t.StrategyId,
		TokensBought:  tokenAmount / constants.TokenAmountDecimals,
		AmountPaid:    solAmount / constants.LamportsConversion,
		TxHash:        transaction.GetSignature(),
		WalletAddress: t.Wallet.PublicKey().String(),
		TokenAddress:  t.Token.String(),
		BondingCurve:  bondingCurve,
	})
	if err != nil {
		logger.Error(err)
	}
}

func notifySell(t *tasks.SellTask, transaction transaction.Transaction, pos *positionModel.Position, tokensSold float64, solReceived float64, deps *Dependencies) {
	bondingCurve, err := pda.GetBondingCurveAddress(t.Token.String())
	if err != nil {
		logger.Error(err)
	}

	strategyTask, err := deps.TradingService.GetBy(*t.StrategyId)
	if err != nil {
		logger.Error(err)
		return
	}

	tokensRemainingRaw, _ := pos.TokenRemaining.Float64()
	currentProfitRaw, _ := pos.FinalizedProfit.Float64()

	err = deps.Notifier.SendSuccessSell(notifierModel.SellNotifierPayload{
		TaskType:        string(strategyTask.StrategyType()),
		TaskId:          t.Id(),
		StrategyId:      t.StrategyId,
		TxHash:          transaction.GetSignature(),
		TokensSold:      tokensSold / constants.TokenAmountDecimals,
		TokensRemaining: tokensRemainingRaw / constants.TokenAmountDecimals,
		CurrentProfit:   currentProfitRaw / constants.LamportsConversion,
		AmountSold:      solReceived / constants.LamportsConversion,
		WalletAddress:   t.Wallet.PublicKey().String(),
		TokenAddress:    t.Token.String(),
		BondingCurve:    bondingCurve,
	})
	if err != nil {
		logger.Error(err)
	}
}

func notifyError(errorMessage string, deps *Dependencies, task tasks.Task, transaction transaction.Transaction) error {
	if task.State().TaskState.IsCancel() {
		return nil
	}

	strategy, err := deps.TradingService.GetBy(*task.GetStrategyId())
	if err != nil {
		return err
	}

	bondingCurve, err := pda.GetBondingCurveAddress(task.GetToken())
	if err != nil {
		return err
	}

	err = deps.Notifier.SendFailure(notifierModel.ErrorNotifierPayload{
		TaskType:      string(strategy.StrategyType()),
		TaskId:        task.Id(),
		StrategyId:    task.GetStrategyId(),
		Error:         errorMessage,
		TxHash:        transaction.GetSignature(),
		WalletAddress: task.GetWallet(),
		TokenAddress:  task.GetToken(),
		BondingCurve:  bondingCurve,
	})

	return err
}
