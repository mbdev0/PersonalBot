package handlers

import (
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/monitoring/models"
)

func HandleCoinFiltering(coin *models.Coin) *models.Coin {
	pipeline := filters.FilterPipeline{}
	// pipeline.AddFilter(filters.HasWebsite())

	result := pipeline.ApplyFilters(coin)
	return result
}
