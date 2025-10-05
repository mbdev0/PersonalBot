package solana_price

import (
	"encoding/json"
	"pump_fun/infrastructure/http"
	"pump_fun/internal/core/constants"
	"pump_fun/pkg/logger"
	"sync"
	"time"
)

var (
	lock             = &sync.Mutex{}
	solPrice         *float64
	priceRefreshTime = 30 * time.Minute
)

func GetSolPrice() (*float64, error) {
	if solPrice == nil {
		lock.Lock()
		defer lock.Unlock()

		if solPrice == nil {
			endPointResponse, err := fetchSolPriceFromEndpoint()
			if err != nil {
				return nil, err
			}

			solPrice = endPointResponse
			go startBackgroundUpdate()
		}
	}

	return solPrice, nil
}

func fetchSolPriceFromEndpoint() (*float64, error) {
	resp, err := http.Get(constants.PumpFunAPIEndPoint + "sol-price")
	if err != nil {
		return nil, err
	}

	var response SolPriceResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		logger.Error("Error unmarshaling JSON", err)
		return nil, err
	}

	return &response.SolPrice, nil
}

func startBackgroundUpdate() {
	ticker := time.NewTicker(priceRefreshTime)
	defer ticker.Stop()

	for range ticker.C {
		lock.Lock()
		epResponse, err := fetchSolPriceFromEndpoint()
		if err != nil {
			return
		}

		solPrice = epResponse

		lock.Unlock()
	}
}
