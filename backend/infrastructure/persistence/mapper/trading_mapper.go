package mapper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/wallets"
	"personal_bot/internal/solana/monitoring/filters"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapTradingTaskToRepo(task strategies.Task) (*models.TradingRow, error) {
	switch t := task.(type) {
	case *strategies.Afk:
		return mapAfkToRepo(t)
	case *strategies.Buy:
		return mapBuyToRepo(t)
	case *strategies.Sell:
		return mapSellToRepo(t)
	default:
		return nil, fmt.Errorf("unknown task type")
	}

}

func mapAfkToRepo(afk *strategies.Afk) (*models.TradingRow, error) {
	taskRow := models.TradingRow{}
	taskRow.ComputeUnits = int(afk.ComputeUnits)
	taskRow.Slippage = int(afk.Slippage * 100)
	taskRow.TradingType = "AFK"
	taskRow.WalletId = afk.Wallet.Id
	taskRow.Id = int(afk.StrategyTaskId())
	taskRow.TimeCreatedUnix = afk.TimeCreated

	afkConfig, err := mapToAfkConfig(afk)
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(afkConfig)
	if err != nil {
		return nil, err
	}

	taskRow.Config = string(json)
	taskRow.RpcGroupId = afk.RPCGroup.Id

	return &taskRow, nil
}

func mapToAfkConfig(afk *strategies.Afk) (models.AfkConfig, error) {
	afkConfig := models.AfkConfig{}
	afkConfig.Filters = mapFiltersToRepo(afk.Filters)

	afkConfig.BuyAmount = int(afk.BuyAmount.Int64())

	buyFee := utils.ConvertSolToLamport(afk.BuyFee)
	afkConfig.BuyFee = int(buyFee.Int64())

	if afk.SellFee != nil {
		sellFee := utils.ConvertSolToLamport(*afk.SellFee)
		sellFeeInt := int(sellFee.Int64())
		afkConfig.SellFee = &sellFeeInt
	}

	afkConfig.SellStrategies = mapSellStratsToRepo(afk.SellStrategies)

	return afkConfig, nil
}

func mapBuyToRepo(buy *strategies.Buy) (*models.TradingRow, error) {
	taskRow := models.TradingRow{}
	taskRow.ComputeUnits = int(buy.ComputeUnits)
	taskRow.Slippage = int(buy.Slippage * 100)
	taskRow.TradingType = "BUY"
	taskRow.WalletId = buy.Wallet.Id
	taskRow.Id = int(buy.StrategyTaskId())
	taskRow.TimeCreatedUnix = buy.TimeCreated

	buyConfig, err := mapToBuyConfig(buy)
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(buyConfig)
	if err != nil {
		return nil, err
	}

	taskRow.Config = string(json)
	taskRow.RpcGroupId = buy.RPCGroup.Id

	return &taskRow, nil

}

func mapToBuyConfig(buy *strategies.Buy) (models.BuyStrategyConfig, error) {
	buyConfig := models.BuyStrategyConfig{}

	buyConfig.BuyAmount = int(buy.BuyAmount.Int64())

	buyFee := utils.ConvertSolToLamport(buy.BuyFee)
	buyConfig.BuyFee = int(buyFee.Int64())

	buyConfig.Token = buy.Token.String()

	if buy.SellFee != nil {
		sellFee := utils.ConvertSolToLamport(*buy.SellFee)
		sellFeeInt := int(sellFee.Int64())
		buyConfig.SellFee = &sellFeeInt
	}

	buyConfig.SellStrategies = mapSellStratsToRepo(buy.SellStrategies)

	buyConfig.BuyTaskId = int(buy.BuyTaskId)
	buyConfig.PositionId = int(buy.PositionId)

	return buyConfig, nil
}

func mapSellToRepo(t *strategies.Sell) (*models.TradingRow, error) {
	taskRow := models.TradingRow{}
	taskRow.ComputeUnits = int(t.ComputeUnits)
	taskRow.Slippage = int(t.Slippage * 100)
	taskRow.TradingType = "SELL"
	taskRow.WalletId = t.Wallet.Id
	taskRow.Id = int(t.StrategyTaskId())
	taskRow.TimeCreatedUnix = t.TimeCreated

	sellConfig, err := mapToSellConfig(t)
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(sellConfig)
	if err != nil {
		return nil, err
	}

	taskRow.Config = string(json)
	taskRow.RpcGroupId = t.RPCGroup.Id

	return &taskRow, nil
}

func mapToSellConfig(t *strategies.Sell) (models.SellStrategyConfig, error) {
	sellConfig := models.SellStrategyConfig{}

	sellConfig.SellAmount = t.SellAmount
	sellFee := utils.ConvertSolToLamport(t.SellFee)
	sellConfig.SellFee = int(sellFee.Int64())

	sellConfig.Token = t.Token.String()
	sellConfig.SellTaskId = int(t.SellTaskId)

	return sellConfig, nil
}

func mapFiltersToRepo(filters []strategies.StrategyFilter) models.Filters {
	dest := models.Filters{}
	for _, filterFunc := range filters {
		filter := filterFunc()
		b := true

		switch filter.Name {
		case "HasWebsite":
			dest.HasWebsite = &b
		case "HasTwitter":
			dest.HasTwitter = &b
		case "HasTelegram":
			dest.HasTelegram = &b
		case "DevWallet":
			if filter.Value != "" {
				dest.DevWallet = &filter.Value
			}
		}
	}

	return dest
}

func mapSellStratsToRepo(src []strategies.StrategyConfig) []models.SellStrategies {
	dest := make([]models.SellStrategies, len(src))

	for i, config := range src {
		dest[i] = models.SellStrategies{
			Type:       string(config.Type),
			Value:      config.Value,
			SellAmount: config.SellAmount,
		}
	}

	return dest
}

func MapRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository) (strategies.Task, error) {
	switch src.TradingType {
	case "AFK":
		return mapAfkRepoToTradingTask(src, wallet, rpcGroup)
	case "BUY":
		return mapBuyRepoToTradingTask(src, wallet, rpcGroup)
	case "SELL":
		return mapSellRepoToTradingTask(src, wallet, rpcGroup)
	default:
		return nil, fmt.Errorf("unknown trading type: %s", src.TradingType)
	}
}

func mapAfkRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository) (*strategies.Afk, error) {
	afkTask := strategies.Afk{}
	afkTask.New()
	afkTask.SetId(int64(src.Id))

	config := models.AfkConfig{}
	err := json.Unmarshal([]byte(src.Config), &config)
	if err != nil {
		return nil, err
	}

	afkTask.State = string(strategies.CREATED)
	afkTask.BuyAmount = big.NewInt(int64(config.BuyAmount))
	afkTask.BuyFee = utils.ConvertLamportToSol(big.NewInt(int64(config.BuyFee)))
	afkTask.ComputeUnits = float64(src.ComputeUnits)
	afkTask.Slippage = float64(src.Slippage) / 100

	privateKey, err := solana.PrivateKeyFromBase58(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	afkTask.Wallet = wallets.SolanaWallet{
		Id:         wallet.Id,
		WalletName: wallet.WalletName,
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}

	afkTask.Filters = mapRepoFiltersToFilters(config)

	if config.SellFee != nil {
		sellFee := utils.ConvertLamportToSol(big.NewInt(int64(*config.SellFee)))
		afkTask.SellFee = &sellFee
	}

	afkTask.SellStrategies = mapRepoToStrategyConfigs(config.SellStrategies)
	afkTask.TimeCreated = src.TimeCreatedUnix

	afkTask.RPCGroup, err = MapRepositoryToRpcGroup(rpcGroup)
	if err != nil {
		return nil, err
	}

	return &afkTask, nil
}

func mapBuyRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository) (*strategies.Buy, error) {
	buyTask := strategies.Buy{}
	buyTask.New()
	buyTask.SetId(int64(src.Id))

	config := models.BuyStrategyConfig{}
	err := json.Unmarshal([]byte(src.Config), &config)
	if err != nil {
		return nil, err
	}

	buyTask.State = string(strategies.CREATED)
	buyTask.BuyAmount = big.NewInt(int64(config.BuyAmount))
	buyTask.BuyFee = utils.ConvertLamportToSol(big.NewInt(int64(config.BuyFee)))
	buyTask.ComputeUnits = float64(src.ComputeUnits)
	buyTask.Slippage = float64(src.Slippage) / 100

	privateKey, err := solana.PrivateKeyFromBase58(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	buyTask.Wallet = wallets.SolanaWallet{
		Id:         wallet.Id,
		WalletName: wallet.WalletName,
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}

	if config.SellFee != nil {
		sellFee := utils.ConvertLamportToSol(big.NewInt(int64(*config.SellFee)))
		buyTask.SellFee = &sellFee
	}
	buyTask.SellStrategies = mapRepoToStrategyConfigs(config.SellStrategies)
	token, err := solana.PublicKeyFromBase58(config.Token)
	if err != nil {
		return nil, err
	}

	buyTask.Token = token
	buyTask.PositionId = int64(config.PositionId)
	buyTask.BuyTaskId = int64(config.BuyTaskId)
	buyTask.TimeCreated = src.TimeCreatedUnix

	buyTask.RPCGroup, err = MapRepositoryToRpcGroup(rpcGroup)
	if err != nil {
		return nil, err
	}

	return &buyTask, nil
}

func mapSellRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository) (strategies.Task, error) {
	sellTask := strategies.Sell{}
	sellTask.New()
	sellTask.SetId(int64(src.Id))

	config := models.SellStrategyConfig{}
	err := json.Unmarshal([]byte(src.Config), &config)
	if err != nil {
		return nil, err
	}

	sellTask.State = string(strategies.CREATED)
	sellTask.SellAmount = config.SellAmount
	sellTask.SellFee = utils.ConvertLamportToSol(big.NewInt(int64(config.SellFee)))
	sellTask.ComputeUnits = float64(src.ComputeUnits)
	sellTask.Slippage = float64(src.Slippage) / 100
	sellTask.SellTaskId = int64(config.SellTaskId)
	sellTask.TimeCreated = src.TimeCreatedUnix

	privateKey, err := solana.PrivateKeyFromBase58(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	sellTask.Wallet = wallets.SolanaWallet{
		Id:         wallet.Id,
		WalletName: wallet.WalletName,
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}

	token, err := solana.PublicKeyFromBase58(config.Token)
	if err != nil {
		return nil, err
	}

	sellTask.Token = token
	sellTask.RPCGroup, err = MapRepositoryToRpcGroup(rpcGroup)
	if err != nil {
		return nil, err
	}

	return &sellTask, nil
}

func mapRepoFiltersToFilters(config models.AfkConfig) []strategies.StrategyFilter {

	srcFilters := config.Filters
	destFilters := make([]strategies.StrategyFilter, 0)
	extractAndCheck := func(f *bool) bool {
		return f != nil && *f
	}

	if extractAndCheck(srcFilters.HasTelegram) {
		destFilters = append(destFilters, filters.HasTelegram)
	}
	if extractAndCheck(srcFilters.HasTwitter) {
		destFilters = append(destFilters, filters.HasTwitter)
	}
	if extractAndCheck(srcFilters.HasWebsite) {
		destFilters = append(destFilters, filters.HasWebsite)
	}

	if srcFilters.DevWallet != nil && *srcFilters.DevWallet != "" {
		wallet := *srcFilters.DevWallet
		destFilters = append(destFilters, func() filters.FilterInfo {
			return filters.DevWallet(wallet)
		})
	}

	return destFilters
}

func mapRepoToStrategyConfigs(sellStrats []models.SellStrategies) []strategies.StrategyConfig {
	configs := make([]strategies.StrategyConfig, len(sellStrats))
	for i, ss := range sellStrats {
		configs[i] = strategies.StrategyConfig{
			Type:       strategies.SellStrategyType(ss.Type),
			Value:      ss.Value,
			SellAmount: ss.SellAmount,
		}
	}
	return configs
}
