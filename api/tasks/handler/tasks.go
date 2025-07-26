package handler

import (
	"encoding/json"
	"net/http"
	"pump_fun/api/models"
	"pump_fun/api/tasks/controller"
)

type TaskHandler struct {
	controller *controller.TaskController
}

func NewTaskHandler(controller *controller.TaskController) http.Handler {
	mux := http.NewServeMux()
	taskHandler := &TaskHandler{controller: controller}
	taskHandler.registerRoutes(mux)
	return mux
}

func (th *TaskHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /create", th.CreateTask)
	mux.HandleFunc("/test", th.Test)
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reqTask models.RequestTask
	err := decoder.Decode(&reqTask)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errString := "Error during JSON decoding: " + err.Error()
		w.Write([]byte(errString))
		return
	}

	createdTask, err := th.controller.CreateTask(reqTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errString := "Error during creation of task: " + err.Error()
		w.Write([]byte(errString))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (th *TaskHandler) Test(w http.ResponseWriter, r *http.Request) {
	res := th.controller.TestEP()
	w.Write([]byte(res))
}
