package handler

import (
	"encoding/json"
	"net/http"
	"pump_fun/api/controller"
	"pump_fun/api/dto"
	"pump_fun/pkg/logger"
)

type WalletsHandler struct {
	controller *controller.WalletsController
}

func NewWalletHandler(ctrl *controller.WalletsController) http.Handler {
	mux := http.NewServeMux()
	walletHandler := WalletsHandler{controller: ctrl}
	walletHandler.registerRoutes(mux)

	return mux
}

func (wh *WalletsHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /wallets/get", wh.getWallets)
	mux.HandleFunc("GET /wallets/get/{id}", wh.getWalletById)
	mux.HandleFunc("POST /wallets/add", wh.addWallet)
	mux.HandleFunc("DELETE /wallets/delete", wh.deleteWallets)
}

func (wh *WalletsHandler) getWallets(w http.ResponseWriter, r *http.Request) {
	wallets, err := wh.controller.GetWallets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(wallets)
	if err != nil {
		logger.Error("Error in getWallets e.p", err)
	}
}

func (wh *WalletsHandler) getWalletById(w http.ResponseWriter, r *http.Request) {
}

func (wh *WalletsHandler) addWallet(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var walletDto dto.RequestWalletDto

	err := decoder.Decode(&walletDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	hasSucceeded, err := wh.controller.InsertWallet(r.Context(), walletDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if hasSucceeded {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
	} else {
		http.Error(w, "Wallet insertion failed", http.StatusInternalServerError)
	}

}

func (wh *WalletsHandler) deleteWallets(w http.ResponseWriter, r *http.Request) {
}
