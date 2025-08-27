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

func createBuyTask(req *dto.RequestTask) (task *tasks.BuyTask, err error) {
	if req.BuyAmount == nil {
		return nil, fmt.Errorf("buy amount not filled in")
	}

	if req.BuyFee == nil {
		return nil, fmt.Errorf("buy fee not filled in")
	}

	bigBuyAmount := big.NewInt(int64(*req.BuyAmount * constants.LamportsConversion))

	wallet, err := solana.PrivateKeyFromBase58(req.WalletAddressPrivateKey)
	if err != nil {
		return nil, err
	}

	address, err := solana.PublicKeyFromBase58(req.TokenAddress)
	if err != nil {
		return nil, err
	}

	bt := tasks.NewBuyTask(wallet, address,
		[]tasks.Option{
			tasks.WithComputeUnits(req.ComputeUnits),
			tasks.WithSlippage(req.Slippage),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(bigBuyAmount),
			tasks.WithBuyFee(*req.BuyFee),
		},
	)

	return bt, nil
}

func createSellTask(reqTask *dto.RequestTask) (task *tasks.SellTask, err error) {
	if reqTask.SellAmount == nil {
		return nil, fmt.Errorf("sell amount is empty")
	}

	if reqTask.SellFee == nil {
		return nil, fmt.Errorf("sell fee is empty")
	}

	token, err := solana.PublicKeyFromBase58(reqTask.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("token address is invalid format")
	}
	wallet, err := solana.PrivateKeyFromBase58(reqTask.WalletAddressPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("private key is not valid/not found")
	}

	sellTask := tasks.NewSellTask(wallet, token,
		[]tasks.Option{
			tasks.WithComputeUnits(reqTask.ComputeUnits),
			tasks.WithSlippage(reqTask.Slippage),
		},
		[]tasks.SellOption{
			tasks.WithSellAmount(*reqTask.SellAmount),
			tasks.WithSellFee(*reqTask.SellFee),
		},
	)

	return sellTask, nil
}
