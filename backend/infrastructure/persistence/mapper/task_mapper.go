package mapper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/utils"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

func MapRepoToTask(task models.TaskRow, wallet models.WalletRepository) (tasks.Task, error) {
	switch task.TaskType {
	case "BUY":
		return createBuyTask(task, wallet)
	case "SELL":
		return createSellTask(task, wallet)
	default:
		return nil, fmt.Errorf("invalid task type in db: %s", task.TaskType)
	}
}

func createBuyTask(src models.TaskRow, wallet models.WalletRepository) (*tasks.BuyTask, error) {
	buyConfig := new(models.BuyConfig)
	if err := json.Unmarshal([]byte(src.Config), buyConfig); err != nil {
		return nil, err
	}

	nodeConfig := new(models.NodeConfig)
	if err := json.Unmarshal([]byte(src.NodeConfig), nodeConfig); err != nil {
		return nil, err
	}

	mappedWallet, err := WalletRepoToWallet(wallet)
	if err != nil {
		return nil, err
	}

	token, err := solana.PublicKeyFromBase58(src.Token)
	if err != nil {
		return nil, err
	}

	if buyConfig.BuyAmount <= 0 {
		return nil, fmt.Errorf("invalid buy amount: %d", buyConfig.BuyAmount)
	}

	buyAmnt := new(big.Int).SetInt64(int64(buyConfig.BuyAmount))
	buyFeeLamport := new(big.Int).SetInt64(int64(buyConfig.BuyFee))
	buyFee := utils.ConvertLamportToSol(buyFeeLamport)

	logger.Information(buyFeeLamport)
	logger.Information(buyFee)

	buyOpts := []tasks.BuyOption{
		tasks.WithBuyAmount(buyAmnt),
		tasks.WithBuyFee(buyFee),
	}

	taskOpts := []tasks.Option{
		tasks.WithSlippage(float64(src.Slippage) / 100.0),
		tasks.WithComputeUnits(uint32(src.ComputeUnits)),
		tasks.WithUnixTime(src.TimeCreatedUnix),
		tasks.WithHttpNode(nodeConfig.Http),
		tasks.WithWS(nodeConfig.Ws),
	}

	if src.StrategyId != nil {
		taskOpts = append(taskOpts, tasks.WithStrategyId(*src.StrategyId))
	}

	buyTask := tasks.NewBuyTask(mappedWallet, token, taskOpts, buyOpts)

	buyTask.SetId(int64(src.Id))

	return buyTask, nil
}

// THIS SHOULD NOT BE CONNECTED TO ANY POSITION
func createSellTask(src models.TaskRow, wallet models.WalletRepository) (*tasks.SellTask, error) {
	sellConfig := new(models.SellConfig)
	if err := json.Unmarshal([]byte(src.Config), sellConfig); err != nil {
		return nil, err
	}

	nodeConfig := new(models.NodeConfig)
	if err := json.Unmarshal([]byte(src.NodeConfig), nodeConfig); err != nil {
		return nil, err
	}

	mappedWallet, err := WalletRepoToWallet(wallet)
	if err != nil {
		return nil, err
	}

	token, err := solana.PublicKeyFromBase58(src.Token)
	if err != nil {
		return nil, err
	}

	if sellConfig.SellAmount <= 0 || sellConfig.SellAmount > 100 {
		return nil, fmt.Errorf("invalid sell amount: %d (must be 1-100)", sellConfig.SellAmount)
	}

	taskOpts := []tasks.Option{
		tasks.WithComputeUnits(uint32(src.ComputeUnits)),
		tasks.WithSlippage(float64(src.Slippage) / 100.0),
		tasks.WithUnixTime(src.TimeCreatedUnix),
		tasks.WithHttpNode(nodeConfig.Http),
		tasks.WithWS(nodeConfig.Ws),
	}

	if src.StrategyId != nil {
		taskOpts = append(taskOpts, tasks.WithStrategyId(*src.StrategyId))
	}

	sellFeeLamport := new(big.Int).SetInt64(int64(sellConfig.SellFee))
	sellFee := utils.ConvertLamportToSol(sellFeeLamport)
	sellTask := tasks.NewSellTask(mappedWallet, token, taskOpts, []tasks.SellOption{
		tasks.WithSellAmount(float64(sellConfig.SellAmount) / 100.0),
		tasks.WithSellFee(sellFee),
	})

	sellTask.SetId(int64(src.Id))

	return sellTask, nil
}

func MapTaskToRepo(src tasks.Task) (*models.TaskRow, error) {
	taskRow := &models.TaskRow{
		Id:    int(src.Id()),
		State: tasks.TaskCreate.ToString(),
	}

	switch task := src.(type) {
	case *tasks.BuyTask:

		if !task.BuyAmount.IsInt64() {
			return nil, fmt.Errorf("buy amount too large to store: %s", task.BuyAmount.String())
		}

		buyFee := utils.ConvertSolToLamport(task.Fee)
		buyConfig := models.BuyConfig{
			BuyFee:    int(buyFee.Int64()),
			BuyAmount: int(task.BuyAmount.Int64()),
		}

		configJSON, err := json.Marshal(buyConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal buy config: %w", err)
		}

		taskRow.TaskType = "BUY"
		taskRow.WalletId = task.WalletId
		taskRow.Token = task.Token.String()
		taskRow.Slippage = int(task.Slippage * 100)
		taskRow.ComputeUnits = int(task.ComputeUnits)
		taskRow.Config = string(configJSON)
		taskRow.StrategyId = task.StrategyId
		taskRow.TimeCreatedUnix = task.TimeCreated

		nodeConfig := models.NodeConfig{
			Http: task.HttpNode(),
			Ws:   task.WSNode(),
		}
		nodeConfigJson, err := json.Marshal(nodeConfig)
		if err != nil {
			return nil, err
		}

		taskRow.NodeConfig = string(nodeConfigJson)

	case *tasks.SellTask:

		sellFeeBig := utils.ConvertSolToLamport(task.Fee)

		sellConfig := models.SellConfig{
			SellFee:    int(sellFeeBig.Int64()),
			SellAmount: int(task.SellPercentage * 100),
		}

		configJSON, err := json.Marshal(sellConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal sell config: %w", err)
		}

		taskRow.TaskType = "SELL"
		taskRow.WalletId = task.WalletID
		taskRow.Token = task.Token.String()
		taskRow.Slippage = int(task.Slippage * 100)
		taskRow.ComputeUnits = int(task.ComputeUnits)
		taskRow.TimeCreatedUnix = task.TimeCreated
		taskRow.Config = string(configJSON)
		taskRow.StrategyId = task.StrategyId

		nodeConfig := models.NodeConfig{
			Http: task.HttpNode(),
			Ws:   task.WSNode(),
		}
		nodeConfigJson, err := json.Marshal(nodeConfig)
		if err != nil {
			return nil, err
		}

		taskRow.NodeConfig = string(nodeConfigJson)

	default:
		return nil, fmt.Errorf("unknown task type: %T", src)
	}

	return taskRow, nil
}
