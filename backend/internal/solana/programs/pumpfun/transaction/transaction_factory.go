package transaction

import (
	"context"
	"fmt"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	ammBuy "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/transaction/buy"
	ammSell "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/transaction/sell"

	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction/buy"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction/sell"

	"personal_bot/internal/solana/transaction"
)

func GetTransactionType(ctx context.Context, task tasks.Task, posService *position.Service, publisher subscriptionhub.Publisher) (transaction.Transaction, error) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		return getBuyImplementation(ctx, t, posService, publisher)
	case *tasks.SellTask:
		return getSellImplementation(ctx, t, posService, publisher)
	}
	return nil, fmt.Errorf("no transaction found for task: %s", task.Type())
}

func getBuyImplementation(ctx context.Context, task *tasks.BuyTask, posService *position.Service, publisher subscriptionhub.Publisher) (transaction transaction.Transaction, err error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(task.GetToken(), constants.PumpFunProgram)
	if err != nil {
		return nil, err
	}

	_, err, hasCompleted := bondingcurve.GetBondingCurveDataFromAddress(ctx, bondingCurveAddress, task.HttpClient())
	if err != nil {
		if hasCompleted {
			return &ammBuy.Transaction{BuyTask: task, PositionService: posService, Publisher: publisher}, nil
		}
		return nil, err
	}
	return &buy.Transaction{BuyTask: task, PositionService: posService, Publisher: publisher}, nil
}

func getSellImplementation(ctx context.Context, task *tasks.SellTask, posService *position.Service, publisher subscriptionhub.Publisher) (transaction transaction.Transaction, err error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(task.GetToken(), constants.PumpFunProgram)
	if err != nil {
		return nil, err
	}

	_, err, hasCompleted := bondingcurve.GetBondingCurveDataFromAddress(ctx, bondingCurveAddress, task.HttpClient())
	if err != nil {
		if hasCompleted {
			return &ammSell.Transaction{SellTask: task, PositionService: posService, Publisher: publisher}, nil
		}
		return nil, err
	}
	return &sell.Transaction{Task: task, PositionService: posService, Publisher: publisher}, nil
}
