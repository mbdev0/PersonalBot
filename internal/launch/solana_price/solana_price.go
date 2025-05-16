package solana_price

import (
	"encoding/json"
	"pump_fun/internal/constants"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	requestclient "pump_fun/internal/request_client"
	"sync"
	"time"
)

var (
	lock     = &sync.Mutex{}
	solPrice *float64
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
	resp, err := requestclient.Get(constants.PumpFunAPIEndPoint + "sol-price")
	if err != nil {
		return nil, err
	}

	var response models.SolPriceResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		logger.Log(logger.LevelError, "Error unmarshaling JSON", logger.String("error", err.Error()))
		return nil, err
	}

	return &response.SolPrice, nil
}

func startBackgroundUpdate() {
	ticker := time.NewTicker(30 * time.Minute)
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
