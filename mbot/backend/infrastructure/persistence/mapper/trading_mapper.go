package mapper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"personal_bot/backend/infrastructure/persistence/models"
	"personal_bot/backend/internal/core/strategies"
	"personal_bot/backend/internal/core/wallets"
	"personal_bot/backend/internal/solana/monitoring/filters"
	"personal_bot/backend/internal/solana/utils"
	jsonparse "personal_bot/backend/pkg/json_parse"
	"personal_bot/backend/pkg/logger"

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
	case *strategies.Spam:
		return mapSpamToRepo(t)
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

	afkConfig.Retries = afk.Retries
	afkConfig.RetriesDelayMs = afk.RetriesDelayMS

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
	buyConfig.Retries = buy.Retries
	buyConfig.RetriesDelayMs = buy.RetriesDelayMS

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

	sellConfig.Retries = t.Retries
	sellConfig.RetriesDelayMs = t.RetriesDelayMS

	return sellConfig, nil
}

func mapSpamToRepo(spam *strategies.Spam) (*models.TradingRow, error) {
	spamRow := models.TradingRow{}
	spamRow.ComputeUnits = int(spam.ComputeUnits)
	spamRow.Slippage = int(spam.Slippage * 100)
	spamRow.TradingType = "SPAM"
	spamRow.WalletId = spam.Wallet.Id
	spamRow.Id = int(spam.StrategyTaskId())
	spamRow.TimeCreatedUnix = spam.TimeCreated
	spamRow.RpcGroupId = spam.RPCGroupId()

	config := mapToSpamConfig(spam)

	configJson, err := jsonparse.Encode(config)
	if err != nil {
		return nil, err
	}
	spamRow.Config = string(configJson)

	logger.Information(spamRow)

	return &spamRow, nil
}

func mapToSpamConfig(spam *strategies.Spam) models.SpamStrategyConfig {
	config := models.SpamStrategyConfig{}

	config.BuyAmount = int(spam.BuyAmount.Int64())
	buyFee := utils.ConvertSolToLamport(spam.BuyFee)
	config.BuyFee = int(buyFee.Int64())

	config.Token = spam.Token.String()

	if spam.SellFee != nil {
		sellFee := utils.ConvertSolToLamport(*spam.SellFee)
		sellFeeInt := int(sellFee.Int64())
		config.SellFee = &sellFeeInt
	}

	config.NumberOfSubTasks = spam.NumberOfSubTasks
	config.Retries = spam.Retries

	config.RetriesDelayMs = spam.RetriesDelayMS

	config.StartTime = spam.StartTime

	return config

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

func MapRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository, program models.ProgramRepository) (strategies.Task, error) {
	switch src.TradingType {
	case "AFK":
		return mapAfkRepoToTradingTask(src, wallet, rpcGroup, program)
	case "BUY":
		return mapBuyRepoToTradingTask(src, wallet, rpcGroup, program)
	case "SELL":
		return mapSellRepoToTradingTask(src, wallet, rpcGroup, program)
	case "SPAM":
		return mapSpamRepoToTradingTask(src, wallet, rpcGroup, program)
	default:
		return nil, fmt.Errorf("unknown trading type: %s", src.TradingType)
	}
}

func mapAfkRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository, program models.ProgramRepository) (*strategies.Afk, error) {
	afkTask := strategies.Afk{}
	afkTask.New()
	afkTask.SetId(int64(src.Id))

	config := models.AfkConfig{}
	err := json.Unmarshal([]byte(src.Config), &config)
	if err != nil {
		return nil, err
	}

	afkTask.Program = program.Program
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

	afkTask.Retries = config.Retries
	afkTask.RetriesDelayMS = config.RetriesDelayMs

	return &afkTask, nil
}

func mapBuyRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository, program models.ProgramRepository) (*strategies.Buy, error) {
	buyTask := strategies.Buy{}
	buyTask.New()
	buyTask.SetId(int64(src.Id))
	buyTask.Program = program.Program

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

	buyTask.Retries = config.Retries
	buyTask.RetriesDelayMS = config.RetriesDelayMs

	return &buyTask, nil
}

func mapSellRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository, program models.ProgramRepository) (strategies.Task, error) {
	sellTask := strategies.Sell{}
	sellTask.New()
	sellTask.SetId(int64(src.Id))
	sellTask.Program = program.Program

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
	sellTask.Retries = config.Retries
	sellTask.RetriesDelayMS = config.RetriesDelayMs

	return &sellTask, nil
}

func mapSpamRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository, rpcGroup models.RpcGroupRepository, program models.ProgramRepository) (strategies.Task, error) {
	spamTask := strategies.Spam{}
	spamTask.New()
	spamTask.SetId(int64(src.Id))

	config := models.SpamStrategyConfig{}
	err := json.Unmarshal([]byte(src.Config), &config)
	if err != nil {
		return nil, err
	}

	spamTask.Program = program.Program
	spamTask.State = string(strategies.CREATED)
	spamTask.BuyAmount = big.NewInt(int64(config.BuyAmount))
	spamTask.BuyFee = utils.ConvertLamportToSol(big.NewInt(int64(config.BuyFee)))
	spamTask.ComputeUnits = float64(src.ComputeUnits)
	spamTask.Slippage = float64(src.Slippage) / 100
	spamTask.NumberOfSubTasks = config.NumberOfSubTasks
	spamTask.StartTime = config.StartTime
	spamTask.Retries = config.Retries
	spamTask.RetriesDelayMS = config.RetriesDelayMs

	token, err := solana.PublicKeyFromBase58(config.Token)
	if err != nil {
		return nil, err
	}

	spamTask.Token = token

	privateKey, err := solana.PrivateKeyFromBase58(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	spamTask.Wallet = wallets.SolanaWallet{
		Id:         wallet.Id,
		WalletName: wallet.WalletName,
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}

	if config.SellFee != nil {
		sellFee := utils.ConvertLamportToSol(big.NewInt(int64(*config.SellFee)))
		spamTask.SellFee = &sellFee
	}

	spamTask.TimeCreated = src.TimeCreatedUnix

	spamTask.RPCGroup, err = MapRepositoryToRpcGroup(rpcGroup)
	if err != nil {
		return nil, err
	}

	return &spamTask, nil
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
