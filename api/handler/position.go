package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pump_fun/api/controller"
	"pump_fun/pkg/logger"
)

type PositionHandler struct {
	controller *controller.PositionController
}

func NewPositionHandler(controller *controller.PositionController) http.Handler {
	mux := http.NewServeMux()
	positionHandler := &PositionHandler{controller: controller}
	positionHandler.registerRoutes(mux)

	return mux
}

func (ph *PositionHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /position/{id}", ph.getPositionById)
	mux.HandleFunc("GET /positions", ph.getPositions)
}

func (ph *PositionHandler) getPositionById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	position, err := ph.controller.GetBy(id)

	//if not found we will return not found
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//otherwise return task
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(position)
	if err != nil {
		logger.Error("error", err)
	}

}

func (ph *PositionHandler) getPositions(w http.ResponseWriter, r *http.Request) {
	allPositions := ph.controller.GetAll()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(allPositions)
	if err != nil {
		logger.Error("error", err)
	}
}
