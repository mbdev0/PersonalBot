package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/api/dto"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/pkg/logger"
	"strconv"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type TaskHandler struct {
	controller *controller.TaskController
	bufferSize int
}

func NewTaskHandler(controller *controller.TaskController) http.Handler {
	mux := http.NewServeMux()
	taskHandler := &TaskHandler{controller: controller}
	taskHandler.registerRoutes(mux)
	taskHandler.bufferSize = 1000

	return mux
}

func (th *TaskHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /task", th.createTask)
	mux.HandleFunc("GET /test", th.Test)
	mux.HandleFunc("GET /task/{id}", th.getTaskById)
	mux.HandleFunc("GET /task", th.getTasks)
	mux.HandleFunc("PUT /task/{id}", th.updateTask)
	mux.HandleFunc("DELETE /task/{id}", th.deleteTask)
	mux.HandleFunc("POST /transition/{id}", th.transitionTask)
	mux.HandleFunc("/subscribe", th.subscribe)
}

func (th *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var reqTask dto.RequestTask

	err := decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	createdTask, err := th.controller.CreateTask(r.Context(), reqTask)
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

func (th *TaskHandler) getTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, fmt.Errorf("invalid id string passed").Error(), http.StatusBadRequest)
		return
	}

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return

	}

	//we have the id
	//we want to talk to controller to get the task relating to the id
	task, err := th.controller.GetTask(int64(convertedId))

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

func (th *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
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

func (th *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	//we extract id from the path
	id := r.PathValue("id")

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	// we need to get the new task
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var reqTask dto.PatchRequestTask
	err = decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//we pass id into controller + new task -> we should be returned the updated task
	updatedTask, err := th.controller.UpdateTask(r.Context(), int64(convertedId), reqTask)

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

func (th *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	err = th.controller.DeleteTask(int64(convertedId))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) transitionTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var transition dto.RequestTransitionTask
	err := decoder.Decode(&transition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	id := r.PathValue("id")
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	err = th.controller.TransitionTask(int64(convertedId), transition.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	subscribers := make(chan subscriptionhub.Subscription, th.bufferSize)
	defer close(subscribers)

	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	}

	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func(c *websocket.Conn) {
		err := c.CloseNow()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			logger.Error(err.Error())
		}
	}(c)

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
		for sub := range subscribers {
			go readFromChannel(sub)
		}
	}()

	//ws loop
	go func() {
		writeToWs := func(msg tasks.TaskEvent) {

			var resp dto.TaskSubResponse
			resp.TaskEvent = &msg
			resp.Error = msg.State.Error

			err := wsjson.Write(r.Context(), c, resp)
			if err != nil {
				logger.Error("error whilst trying to write to WS, error: " + err.Error())
			}
		}
		for msg := range wsWrite {
			writeToWs(msg)
		}
	}()

	//read loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		var msg *dto.TaskSubscribe
		resp := dto.TaskSubResponse{}
		err := wsjson.Read(ctx, c, &msg)

		if err != nil {
			closeStatus := websocket.CloseStatus(err)
			isValidCloseStatus := closeStatus == websocket.StatusNormalClosure ||
				closeStatus == websocket.StatusGoingAway ||
				closeStatus == websocket.StatusNoStatusRcvd

			if isValidCloseStatus {
				logger.Information(fmt.Sprintf("WebSocket closed normally (status: %s)", closeStatus.String()))
				return
			}

			logger.Error("error whilst reading", err)
			resp.Error = err.Error()
			wsjson.Write(ctx, c, "err")
			err := wsjson.Write(ctx, c, resp)
			if err != nil {
				logger.Error("error whilst trying to write to WS, error: " + err.Error())
			}
			return
		}

		switch msg.Type {
		case dto.Subscribe:
			var sub *subscriptionhub.Subscription
			var err error

			if msg.StrategyId != nil {
				sub, err = th.controller.SubscribeByStrategy(*msg.StrategyId)
			} else if msg.Id != nil {
				sub, err = th.controller.Subscribe(*msg.Id)
			} else {
				err = fmt.Errorf("invalid arguments passed: %v", msg)
			}

			if err != nil {
				resp.Error = err.Error()
				err := wsjson.Write(ctx, c, resp)
				if err != nil {
					logger.Error("error whilst trying to write to WS, error: " + err.Error())
				}
				continue
			}
			subscribers <- *sub

		case dto.Unsubscribe:
			var err error
			if msg.StrategyId != nil {
				err = th.controller.UnsubscribeByStrategy(*msg.StrategyId)
			} else if msg.Id != nil {
				err = th.controller.Unsubcribe(*msg.Id)
			} else {
				err = fmt.Errorf("invalid arguments passed: %v", msg)
			}

			if err != nil {
				resp.Error = err.Error()
				err := wsjson.Write(ctx, c, resp)
				if err != nil {
					logger.Error("error whilst trying to write to WS, error: " + err.Error())
				}
				continue
			}
		default:
			resp.Error = "message type wasn't 'Subscribe' or 'Unsubscribe"
			err := wsjson.Write(ctx, c, resp)
			if err != nil {
				logger.Error("error whilst trying to write to WS, error: " + err.Error())
			}
		}
	}
}

func (th *TaskHandler) Test(w http.ResponseWriter, r *http.Request) {
	res := th.controller.TestEP()
	_, err := w.Write([]byte(res))

	if err != nil {
		logger.Error(err.Error())
	}
}
