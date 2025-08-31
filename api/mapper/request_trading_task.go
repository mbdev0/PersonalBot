package mapper

import (
	"fmt"
	"pump_fun/api/dto"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapTradingTaskDtoToTradingTask(src dto.TradingTask) (strategies.Task, error) {
	switch src.Type {
	case dto.AFK:
		task, err := mapAfkDtoToAfk(src)
		if err != nil {
			return nil, err
		}
		return task, nil
	default:
		return nil, fmt.Errorf("task with type: %s - not found", src.Type)
	}
}

func mapAfkDtoToAfk(src dto.TradingTask) (dst *strategies.Afk, err error) {
	dest := strategies.Afk{}
	dest.New()

	dest.BuyFee = src.BuyFee
	dest.ComputeUnits = float64(src.ComputeUnits)
	dest.Slippage = src.Slippage
	dest.BuyAmount = utils.ConvertSolToLamport(src.BuyAmount)
	dest.Wallet, err = solana.PrivateKeyFromBase58(src.Wallet)
	if err != nil {
		return nil, fmt.Errorf("error whilst mapping wallet %w", err)
	}

	dest.Filters = mapFiltersToDestFilters(src.Filters)

	return &dest, nil
}

func mapFiltersToDestFilters(srcFilters dto.Filters) []strategies.StrategyFilter {
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

	return destFilters
}
