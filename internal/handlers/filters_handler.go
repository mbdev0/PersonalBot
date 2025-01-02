package handlers

import (
	"fmt"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/filters"
)

func HandleCoinFiltering(coin *models.Coin) *models.Coin {
	pipeline := filters.FilterPipeline{}
	pipeline.AddFilter(filters.HasWebsite())

	result := pipeline.ApplyFilters(coin)
	//TODO: Remove for debugging
	if result != nil {
		fmt.Printf("Coin did pass filter: %+v\n", *coin)
	} else {
		fmt.Printf("Coin did not pass filter: %+v\n", *coin)
	}
	return result
}
