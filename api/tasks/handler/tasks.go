package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pump_fun/api/models"
	"pump_fun/api/tasks/controller"
	"pump_fun/pkg/logger"
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
	mux.HandleFunc("GET /test", th.Test)
	mux.HandleFunc("GET /task/{id}", th.GetTaskById)
	mux.HandleFunc("GET /task", th.GetTasks)
	mux.HandleFunc("PUT /task/{id}", th.UpdateTask)
	mux.HandleFunc("DELETE /task/{id}", th.DeleteTask)
	mux.HandleFunc("POST /transition/{id}", th.TransitionTask)
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var reqTask models.RequestTask
	err := decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	createdTask, err := th.controller.CreateTask(reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (th *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	//we have the id
	//we want to talk to controller to get the task relating to the id
	task, err := th.controller.GetTask(id)

	//if not found we will return not found
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//otherwise return task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		logger.Error("error", err)
	}

}

func (th *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	allTasks, err := th.controller.GetAllTasks()

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

func (th *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	//we extract id from the path
	id := r.PathValue("id")

	// we need to get the new task
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var reqTask models.RequestTask
	err := decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//we pass id into controller + new task -> we should be returned the updated task
	updatedTask, err := th.controller.UpdateTask(id, reqTask)

	if err != nil {
		logger.Error("Error whilst updating task with id: " + id + " error: " + err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

func (th *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := th.controller.DeleteTask(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) TransitionTask(w http.ResponseWriter, r *http.Request) {
	// we get the body which would be in the format:
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var transition models.RequestTransitionTask
	err := decoder.Decode(&transition)

	id := r.PathValue("id")

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	th.controller.TransitionTask(id, transition.State)
}

func (th *TaskHandler) Test(w http.ResponseWriter, r *http.Request) {
	res := th.controller.TestEP()
	w.Write([]byte(res))
}
