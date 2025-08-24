package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pump_fun/api/controller"
	"pump_fun/api/dto"
	"pump_fun/internal/core/tasks"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/pkg/logger"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type TaskHandler struct {
	controller  *controller.TaskController
	subscribers chan subscriptionhub.Subscription
	bufferSize  int
}

func NewTaskHandler(controller *controller.TaskController) http.Handler {
	mux := http.NewServeMux()
	taskHandler := &TaskHandler{controller: controller}
	taskHandler.registerRoutes(mux)
	taskHandler.bufferSize = 1000
	taskHandler.subscribers = make(chan subscriptionhub.Subscription, taskHandler.bufferSize)

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
	mux.HandleFunc("/subscribe", th.Subscribe)
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var reqTask dto.RequestTask
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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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
	var reqTask dto.PatchRequestTask
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var transition dto.RequestTransitionTask
	err := decoder.Decode(&transition)

	id := r.PathValue("id")

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = th.controller.TransitionTask(id, transition.State)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Error(err)
}

func (th *TaskHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer c.CloseNow()

	wsWrite := make(chan tasks.TaskEvent, 1000)
	defer close(wsWrite)

	// read from each subs channels and write to wsWrite channel
	// writes are threadsafe as channels are built with sync
	readFromChannel := func(s subscriptionhub.Subscription) {
		for tskEvent := range s.Chan() {
			select {
			case wsWrite <- tskEvent:
			default:
				// drop message that couldn't be pushed to wsWrite
				logger.Information("couldn't upload task event to wsWrite - going to drop the event")
			}
		}

	}

	//fan in loop
	go func() {
		for sub := range th.subscribers {
			go readFromChannel(sub)
		}
	}()

	//ws loop
	go func() {
		writeToWs := func(msg tasks.TaskEvent) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var resp dto.SubResponse
			resp.TaskEvent = &msg

			if err != nil {
				resp.Error = err.Error()
				wsjson.Write(ctx, c, resp)
			}
			wsjson.Write(ctx, c, resp)
		}
		for msg := range wsWrite {
			writeToWs(msg)
		}
	}()

	//read loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		var msg *dto.Subscribe
		resp := dto.SubResponse{}
		err := wsjson.Read(ctx, c, &msg)

		if err != nil {
			resp.Error = err.Error()
			wsjson.Write(ctx, c, resp)
			return
		}

		switch msg.Type {
		case "Subscribe":
			sub, err := th.controller.Subscribe(msg.Id)
			if err != nil {
				resp.Error = err.Error()
				wsjson.Write(ctx, c, resp)
				continue
			}
			th.subscribers <- *sub

		case "Unsubscribe":
			err := th.controller.Unsubcribe(msg.Id)
			if err != nil {
				resp.Error = err.Error()
				wsjson.Write(ctx, c, resp)
				continue
			}
		default:
			resp.Error = "message type wasn't 'Subscribe' or 'Unsubscribe"
			wsjson.Write(ctx, c, resp)
		}
	}

}

func (th *TaskHandler) Test(w http.ResponseWriter, r *http.Request) {
	res := th.controller.TestEP()
	w.Write([]byte(res))
}
