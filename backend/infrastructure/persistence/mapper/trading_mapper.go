package mapper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"pump_fun/infrastructure/persistence/models"
	"pump_fun/internal/core/models/wallets"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapTradingTaskToRepo() {}

func MapRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository) (strategies.Task, error) {
	switch src.TradingType {
	case "AFK":
		return mapAfkRepoToTradingTask(src, wallet)
	default:
		return nil, fmt.Errorf("unknown trading type: %s", src.TradingType)
	}
}

func mapAfkRepoToTradingTask(src models.TradingRow, wallet models.WalletRepository) (*strategies.Afk, error) {
	afkTask := strategies.Afk{}
	config := models.AfkConfig{}
	err := json.Unmarshal([]byte(src.Config), &config)
	if err != nil {
		return nil, err
	}

	afkTask.BuyAmount = big.NewInt(int64(config.BuyAmount))
	afkTask.BuyFee = utils.ConvertLamportToSol(big.NewInt(int64(config.BuyFee)))
	afkTask.ComputeUnits = float64(src.ComputeUnits)
	afkTask.Slippage = float64(src.Slippage / 100)

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
	afkTask.SellFee = utils.ConvertLamportToSol(big.NewInt(int64(config.SellFee)))
	afkTask.SellStrategies = mapRepoToStrategyConfigs(config.SellStrategies)

	return &afkTask, nil
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
