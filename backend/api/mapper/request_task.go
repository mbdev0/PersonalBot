package mapper

import (
	"fmt"
	"math/big"
	"personal_bot/api/dto"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/core/tasks"
	"time"

	"github.com/gagliardetto/solana-go"
)

func MapRequestTaskToTask(reqTask *dto.RequestTask, wallet wallets.SolanaWallet) (tasks.Task, error) {
	switch reqTask.Type {
	case dto.Buy:
		buyTask, err := createBuyTask(reqTask, wallet)
		if err != nil {
			return nil, err
		}
		return buyTask, nil
	case dto.Sell:
		sellTask, err := createSellTask(reqTask, wallet)
		if err != nil {
			return nil, err
		}
		return sellTask, nil
	}

	return nil, fmt.Errorf("type of transaction was wrong")
}

func createBuyTask(req *dto.RequestTask, wallet wallets.SolanaWallet) (task *tasks.BuyTask, err error) {
	if req.BuyAmount == nil {
		return nil, fmt.Errorf("buy amount not filled in")
	}

	if req.BuyFee == nil {
		return nil, fmt.Errorf("buy fee not filled in")
	}

	bigBuyAmount := big.NewInt(int64(*req.BuyAmount * constants.LamportsConversion))

	address, err := solana.PublicKeyFromBase58(req.TokenAddress)
	if err != nil {
		return nil, err
	}

	bt := tasks.NewBuyTask(wallet, address,
		[]tasks.Option{
			tasks.WithComputeUnits(req.ComputeUnits),
			tasks.WithSlippage(req.Slippage),
			tasks.WithUnixTime(time.Now().Unix()),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(bigBuyAmount),
			tasks.WithBuyFee(*req.BuyFee),
		},
	)

	return bt, nil
}

func createSellTask(reqTask *dto.RequestTask, wallet wallets.SolanaWallet) (task *tasks.SellTask, err error) {
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

	sellOptions := []tasks.SellOption{
		tasks.WithSellAmount(*reqTask.SellAmount),
		tasks.WithSellFee(*reqTask.SellFee),
	}

	if reqTask.SellPositionId != nil {
		positionOpt := tasks.WithSellPositionId(reqTask.SellPositionId)
		sellOptions = append(sellOptions, positionOpt)
	}

	sellTask := tasks.NewSellTask(wallet, token,
		[]tasks.Option{
			tasks.WithComputeUnits(reqTask.ComputeUnits),
			tasks.WithSlippage(reqTask.Slippage),
			tasks.WithUnixTime(time.Now().Unix()),
		},
		sellOptions,
	)

	return sellTask, nil
}
