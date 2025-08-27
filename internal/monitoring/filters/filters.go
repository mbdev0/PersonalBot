package filters

import (
	"pump_fun/internal/monitoring/models"
)

func HandleCoinFiltering(coin *models.Coin) *models.Coin {
	pipeline := FilterPipeline{}
	pipeline.AddFilter(HasWebsite())

	result := pipeline.ApplyFilters(coin)
	return result
}
