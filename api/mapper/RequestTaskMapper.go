package mapper

import (
	"fmt"
	"math/big"
	"pump_fun/api/dto"
	"pump_fun/internal/core/constants"
	"pump_fun/internal/core/tasks"

	"github.com/gagliardetto/solana-go"
)

func MapRequestTaskToTask(reqTask *dto.RequestTask) (tasks.Task, error) {
	switch reqTask.Type {
	case dto.Buy:
		buyTask, err := createBuyTask(reqTask)
		if err != nil {
			return nil, err
		}
		return buyTask, nil
	case dto.Sell:
		sellTask, err := createSellTask(reqTask)
		if err != nil {
			return nil, err
		}
		return sellTask, nil
	}

	return nil, fmt.Errorf("type of transaction was wrong")
}

func createBuyTask(reqTask *dto.RequestTask) (task *tasks.BuyTask, err error) {
	buyTask := tasks.BuyTask{}
	buyTask.InitDefaults()

	if reqTask.BuyAmount == nil {
		return nil, fmt.Errorf("buy amount not filled in")
	}

	if reqTask.BuyFee == nil {
		return nil, fmt.Errorf("buy fee not filled in")
	}

	bigBuyAmount := big.NewInt(int64((*reqTask.BuyAmount * constants.LamportsConversion)))
	buyTask.BuyAmount = *bigBuyAmount
	buyTask.BuyFee = *reqTask.BuyFee
	buyTask.ComputeUnits = reqTask.ComputeUnits
	buyTask.Slippage = reqTask.Slippage

	buyTask.Wallet, err = solana.PrivateKeyFromBase58(reqTask.WalletAddressPrivateKey)
	if err != nil {
		return nil, err
	}

	buyTask.TokenAddress, err = solana.PublicKeyFromBase58(reqTask.TokenAddress)
	if err != nil {
		return nil, err
	}

	buyTask.State.TaskState = tasks.TaskCreate

	return &buyTask, nil
}

func createSellTask(reqTask *dto.RequestTask) (task *tasks.SellTask, err error) {
	sellTask := tasks.SellTask{}
	sellTask.InitDefaults()

	if reqTask.SellAmount == nil {
		return nil, fmt.Errorf("sell amount is empty")
	}

	if reqTask.SellFee == nil {
		return nil, fmt.Errorf("sell fee is empty")
	}

	sellTask.SellFee = *reqTask.SellFee
	sellTask.PercentageToSell = *reqTask.SellAmount
	sellTask.ComputeUnits = reqTask.ComputeUnits
	sellTask.Slippage = reqTask.Slippage
	sellTask.TokenAddress, err = solana.PublicKeyFromBase58(reqTask.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("token address is invalid format")
	}
	sellTask.Wallet, err = solana.PrivateKeyFromBase58(reqTask.WalletAddressPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("private key is not valid/not found")
	}

	sellTask.State.TaskState = tasks.TaskCreate

	return &sellTask, nil
}
