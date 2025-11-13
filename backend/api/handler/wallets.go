package handler

import (
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("GET /wallets", wh.getWallets)
	mux.HandleFunc("GET /wallets/getById/{id}", wh.getWalletById)
	mux.HandleFunc("GET /wallets/getByName/{name}", wh.getWalletByName)
	mux.HandleFunc("POST /wallets", wh.addWallet)
	mux.HandleFunc("DELETE /wallets/{id}", wh.deleteWallets)
	mux.HandleFunc("PUT /wallets/{id}", wh.updateWallet)
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
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	wallet, err := wh.controller.GetWalletById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		logger.Error("Error in get walletById", err)
	}

}

func (wh *WalletsHandler) getWalletByName(w http.ResponseWriter, r *http.Request) {
	walletName := r.PathValue("name")
	if walletName == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	wallet, err := wh.controller.GetWalletByName(r.Context(), walletName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		logger.Error("Error in get walletByName", err)
	}
}

func (wh *WalletsHandler) deleteWallets(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	deleted, err := wh.controller.DeleteWallet(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(map[string]bool{"IsDeleted": deleted})
	if err != nil {
		logger.Error("Error in get walletByName", err)
	}

}

func (wh *WalletsHandler) updateWallet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var walletDto dto.RequestWalletDto

	err := decoder.Decode(&walletDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	updatedWallet, err := wh.controller.UpdateWallet(r.Context(), id, walletDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedWallet)
	if err != nil {
		logger.Error(err)
	}
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
