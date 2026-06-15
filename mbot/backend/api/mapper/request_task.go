package mapper

import (
	"fmt"
	"math/big"
	"personal_bot/backend/api/dto"
	"personal_bot/backend/internal/core/constants"
	rpcgroups "personal_bot/backend/internal/core/rpc_groups"
	"personal_bot/backend/internal/core/tasks"
	"personal_bot/backend/internal/core/wallets"
	"time"

	"github.com/gagliardetto/solana-go"
)

func MapRequestTaskToTask(reqTask *dto.RequestTask, wallet wallets.SolanaWallet, rpcGroup rpcgroups.GroupItem) (tasks.ConfigurableTask, error) {
	switch reqTask.Type {
	case dto.Buy:
		buyTask, err := createBuyTask(reqTask, wallet, rpcGroup)
		if err != nil {
			return nil, err
		}
		return buyTask, nil
	case dto.Sell:
		sellTask, err := createSellTask(reqTask, wallet, rpcGroup)
		if err != nil {
			return nil, err
		}
		return sellTask, nil
	}

	return nil, fmt.Errorf("type of transaction was wrong")
}

func createBuyTask(req *dto.RequestTask, wallet wallets.SolanaWallet, rpcGroup rpcgroups.GroupItem) (task *tasks.BuyTask, err error) {
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
			tasks.WithProgram(string(req.Program)),
			tasks.WithComputeUnits(req.ComputeUnits),
			tasks.WithSlippage(req.Slippage),
			tasks.WithUnixTime(time.Now().Unix()),
			tasks.WithRPCGroupId(req.RPCGroupId),
			tasks.WithHttpNode(rpcGroup.Http),
			tasks.WithWS(rpcGroup.WS),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(bigBuyAmount),
			tasks.WithBuyFee(*req.BuyFee),
		},
	)

	return bt, nil
}

func createSellTask(reqTask *dto.RequestTask, wallet wallets.SolanaWallet, rpcGroup rpcgroups.GroupItem) (task *tasks.SellTask, err error) {
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

	taskOpts := []tasks.Option{
		tasks.WithProgram(string(reqTask.Program)),
		tasks.WithComputeUnits(reqTask.ComputeUnits),
		tasks.WithSlippage(reqTask.Slippage),
		tasks.WithUnixTime(time.Now().Unix()),
		tasks.WithRPCGroupId(reqTask.RPCGroupId),
		tasks.WithHttpNode(rpcGroup.Http),
		tasks.WithWS(rpcGroup.WS),
	}

	if reqTask.StrategyId != nil {
		taskOpts = append(taskOpts, tasks.WithStrategyId(*reqTask.StrategyId))
	}

	sellTask := tasks.NewSellTask(wallet, token,
		taskOpts,
		sellOptions,
	)

	return sellTask, nil
}
