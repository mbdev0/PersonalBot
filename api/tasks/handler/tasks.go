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
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

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

func (th *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	//we have the id
	//we want to talk to controller to get the task relating to the id
	task, err := th.controller.GetTask(id)

	//if not found we will return not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error: " + err.Error()))
		return
	}

	//otherwise return task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)

}

func (th *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	allTasks := th.controller.GetAllTasks()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTasks)
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON in bad format - check format again"))
	}

	//we pass id into controller + new task -> we should be returned the updated task
	updatedTask, err := th.controller.UpdateTask(id, reqTask)

	if err != nil {
		logger.Error("Error whilst updating task with id: " + id + " error: " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task update failed: " + id + "\nError" + err.Error()))
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("error whilst deleting task id: %s, error: %s", id, err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) Test(w http.ResponseWriter, r *http.Request) {
	res := th.controller.TestEP()
	w.Write([]byte(res))
}
