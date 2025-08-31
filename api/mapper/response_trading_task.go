package mapper

import (
	"fmt"
	"pump_fun/api/dto"
	"pump_fun/internal/core/strategies"
	"pump_fun/internal/solana/utils"
)

func MapTradingTaskToDto(src strategies.Task) (dest *dto.TradingTaskResponse, err error) {
	switch t := src.(type) {
	case *strategies.Afk:
		task, err := mapAfkToAfkResponse(t)
		if err != nil {
			return nil, err
		}
		return task, nil
	default:
		return nil, fmt.Errorf("task not found with type - make sure the type is created")
	}
}

func mapAfkToAfkResponse(src *strategies.Afk) (dest *dto.TradingTaskResponse, err error) {
	dst := dto.TradingTaskResponse{}
	dst.Type = dto.TradingType(src.StrategyType())
	dst.Id = src.StrategyTaskId()
	dst.BuyAmount = utils.ConvertLamportToSol(src.BuyAmount)
	dst.BuyFee = src.BuyFee
	dst.ComputeUnits = src.ComputeUnits
	dst.Slippage = src.Slippage
	dst.Filters = mapFiltersToResponseFilters(src.Filters)

	return &dst, nil
}

func mapFiltersToResponseFilters(src []strategies.StrategyFilter) dto.Filters {
	dest := dto.Filters{}

	for _, filterFunc := range src {
		filter := filterFunc()
		b := true

		switch filter.Name {
		case "HasWebsite":
			dest.HasWebsite = &b
		case "HasTwitter":
			dest.HasTwitter = &b
		case "HasTelegram":
			dest.HasTelegram = &b
		}
	}

	return dest
}
