package handler

import (
	"encoding/json"
	"net/http"
	"personal_bot/backend/api/controller"
	"personal_bot/backend/api/dto"
	"personal_bot/backend/pkg/logger"
	"strconv"
)

type RPCGroupHandler struct {
	controller *controller.RPCGroupController
}

func NewRPCGroupHandler(controller *controller.RPCGroupController) *http.ServeMux {
	mux := http.NewServeMux()

	handler := &RPCGroupHandler{controller: controller}
	handler.registerRoutes(mux)
	return mux
}

func (rgh *RPCGroupHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", rgh.GetRPCGroupDashboard)
	mux.HandleFunc("POST /", rgh.PostRPCGroup)
	mux.HandleFunc("GET /{id}", rgh.GetRPCGroup)
	mux.HandleFunc("DELETE /{id}", rgh.DeleteRPCGroup)
	mux.HandleFunc("PUT /{id}", rgh.PutRPCGroup)
}

func (rgh *RPCGroupHandler) GetRPCGroupDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard, err := rgh.controller.GetDashboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dashboard)
	if err != nil {
		logger.Error(err)
	}
}

func (rgh *RPCGroupHandler) PostRPCGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var rpcGroup dto.RPCGroupPush
	err := decoder.Decode(&rpcGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdGroup, err := rgh.controller.Post(r.Context(), rpcGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdGroup)
	if err != nil {
		logger.Error(err)
	}
}

func (rgh *RPCGroupHandler) GetRPCGroup(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "no id passed in", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "incorrect format of id, must be an integer", http.StatusBadRequest)
		return

	}

	rpcGroup, err := rgh.controller.GetBy(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rpcGroup)
	if err != nil {
		logger.Error(err)
	}
}

func (rgh *RPCGroupHandler) DeleteRPCGroup(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "no id passed in", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "incorrect format of id, must be an integer", http.StatusBadRequest)
		return

	}

	err = rgh.controller.Delete(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rgh *RPCGroupHandler) PutRPCGroup(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	var rpcGroupPush dto.RPCGroupPush
	err = json.NewDecoder(r.Body).Decode(&rpcGroupPush)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	rcpGroup, err := rgh.controller.Update(r.Context(), int64(id), rpcGroupPush)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rcpGroup)
	if err != nil {
		logger.Error(err)
	}
}
