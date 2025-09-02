package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pump_fun/api/controller"
	"pump_fun/api/dto"
	"pump_fun/pkg/logger"

	"github.com/google/uuid"
)

type TradingHandler struct {
	strategyController *controller.StrategyController
}

func NewTradingHandler(controller *controller.StrategyController) http.Handler {
	mux := http.NewServeMux()
	handler := &TradingHandler{strategyController: controller}
	handler.registerRoutes(mux)
	return mux
}

func (th *TradingHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /create", th.createTask)
	mux.HandleFunc("GET /task/{id}", th.getTaskById)
	mux.HandleFunc("GET /task", th.getTasks)
	mux.HandleFunc("PUT /task/{id}", th.updateTask)
	mux.HandleFunc("DELETE /task/{id}", th.deleteTask)
	mux.HandleFunc("GET /task/start/{id}", th.startTask)
	mux.HandleFunc("GET /task/stop/{id}", th.stopTask)

}

func (th *TradingHandler) createTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var reqTask dto.TradingTask

	err := decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	createdTask, err := th.strategyController.Create(reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(createdTask)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func (th *TradingHandler) getTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	//we have the id
	//we want to talk to strategyController to get the task relating to the id
	task, err := th.strategyController.GetBy(id)

	//if not found we will return not found
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//otherwise return task
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		logger.Error("error", err)
	}

}

func (th *TradingHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	allTasks, err := th.strategyController.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusPartialContent)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(allTasks)
	if err != nil {
		logger.Error("error", err)
	}
}

func (th *TradingHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	//we extract id from the path
	id := r.PathValue("id")

	// we need to get the new task
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var reqTask dto.TradingTaskPatch
	err := decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//we pass id into strategyController + new task -> we should be returned the updated task
	updatedTask, err := th.strategyController.Update(id, reqTask)

	if err != nil {
		logger.Error("Error whilst updating task with id: " + id + " error: " + err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTask)
	if err != nil {
		logger.Error("error whilst encoding task", err)
		return
	}
}

func (th *TradingHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := th.strategyController.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TradingHandler) startTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)
	if err != nil {
		http.Error(w, "parameter value for id is not a valid uuid", http.StatusBadRequest)
		return
	}

	err = th.strategyController.Start(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error whilst starting task: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TradingHandler) stopTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)
	if err != nil {
		http.Error(w, "parameter value for id is not a valid uuid", http.StatusBadRequest)
		return
	}

	err = th.strategyController.Stop(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error whilst stopping task: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
