package solana_price

import (
	"encoding/json"
	"personal_bot/infrastructure/http"
	"personal_bot/internal/core/constants"
	"personal_bot/pkg/logger"
	"sync"
	"time"
)

var (
	lock             = &sync.Mutex{}
	solPrice         *float64
	priceRefreshTime = 30 * time.Second
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
	resp, err := http.Get(constants.RadyiumApiEndPoint + "price")
	if err != nil {
		return nil, err
	}

	var response SolPriceResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		logger.Error("Error unmarshaling JSON", err)
		return nil, err
	}

	price, ok := response[constants.SolanaTokenAddress]
	if !ok {
		logger.Error("error whilst getting solana price")
	}

	// logger.Information(fmt.Sprintf("SOLANA PRICE: %f", price))

	return &price, nil
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
