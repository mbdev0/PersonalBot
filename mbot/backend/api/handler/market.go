package handler

import (
	"encoding/json"
	"net/http"
	"personal_bot/backend/api/dto"
	"personal_bot/backend/infrastructure/solana_price"
	"personal_bot/backend/pkg/logger"
)

type MarketHandler struct {
}

func NewMarketHandler() http.Handler {
	mux := http.NewServeMux()
	marketHandler := MarketHandler{}
	marketHandler.registerRoutes(mux)
	return mux
}

func (mh *MarketHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /sol-price", mh.getSolPrice)
}

func (mh *MarketHandler) getSolPrice(w http.ResponseWriter, r *http.Request) {

	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if solPrice == nil {
		http.Error(w, "sol price was nil", http.StatusInternalServerError)
		return
	}

	priceStruct := dto.Price{SolanaPrice: *solPrice}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(priceStruct)
	if err != nil {
		logger.Error(err)
	}
}
