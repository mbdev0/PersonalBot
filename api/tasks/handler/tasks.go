package handler

import (
	"encoding/json"
	"net/http"
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
	mux.HandleFunc("/get", th.GetTasks)
	mux.HandleFunc("/test", th.Test)
}

func (th *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := th.controller.GetTasks()
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (th *TaskHandler) Test(w http.ResponseWriter, r *http.Request) {
	res := th.controller.TestEP()
	w.Write([]byte(res))
}
