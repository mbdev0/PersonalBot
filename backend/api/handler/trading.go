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
	"personal_bot/api/mapper"
	"personal_bot/internal/services/subscription_hub/strategy"
	"personal_bot/pkg/logger"
	"strconv"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type TradingHandler struct {
	strategyController *controller.StrategyController
	bufferSize         int
}

func NewTradingHandler(controller *controller.StrategyController) http.Handler {
	mux := http.NewServeMux()
	handler := &TradingHandler{strategyController: controller}
	handler.bufferSize = 1000
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
	mux.HandleFunc("GET /subscribe", th.subscribe)

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

	createdTask, err := th.strategyController.Create(r.Context(), reqTask)
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
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	//we have the id
	//we want to talk to strategyController to get the task relating to the id
	task, err := th.strategyController.GetBy(int64(convertedId))

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
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	// we need to get the new task
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var reqTask dto.TradingTask
	err = decoder.Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//we pass id into strategyController + new task -> we should be returned the updated task
	updatedTask, err := th.strategyController.Update(r.Context(), int64(convertedId), reqTask)

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
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}
	err = th.strategyController.Delete(r.Context(), int64(convertedId))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TradingHandler) startTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	err = th.strategyController.Start(int64(convertedId))
	if err != nil {
		http.Error(w, fmt.Sprintf("error whilst starting task: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TradingHandler) stopTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id passed").Error(), http.StatusBadRequest)
		return
	}

	err = th.strategyController.Stop(int64(convertedId))
	if err != nil {
		http.Error(w, fmt.Sprintf("error whilst stopping task: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TradingHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	subscribers := make(chan strategy.Subscription, 1000)
	activeSubs := map[int64]strategy.Subscription{}

	defer close(subscribers)
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	}

	//upgrade to ws
	c, err := websocket.Accept(w, r, opts)

	if err != nil {
		http.Error(w, "error whilst trying to transition to WS", http.StatusInternalServerError)
		return
	}

	defer func(c *websocket.Conn) {
		err := c.CloseNow()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			logger.Error(err.Error())
		}
	}(c)

	defer func() {
		for strategyId := range activeSubs {
			err := th.strategyController.Unsubscribe(strategyId)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wsWriteChan := make(chan dto.StrategyResponse, 1000)
	defer close(wsWriteChan)

	//fan in loop
	go th.fanIn(subscribers, wsWriteChan)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	//ws write loop
	go th.writeToWs(ctx, wsWriteChan, c)

	for {
		var msg dto.StrategySubscribe
		resp := dto.StrategyResponse{}

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
			err := wsjson.Write(ctx, c, resp)
			if err != nil {
				logger.Error("error whilst trying to write to WS, error: " + err.Error())
			}
			return
		}

		switch msg.Type {
		case dto.Subscribe:
			sub, err := th.strategyController.Subscribe(msg.Id)
			if err != nil {
				th.handleError(ctx, err, resp, c)
				continue
			}
			subscribers <- *sub
			activeSubs[sub.Sub_id] = *sub

		case dto.Unsubscribe:
			err := th.strategyController.Unsubscribe(msg.Id)
			if err != nil {
				th.handleError(ctx, err, resp, c)
				continue
			}
			delete(activeSubs, msg.Id)
		}
	}

}

func (th *TradingHandler) fanIn(subscribers <-chan strategy.Subscription, wsWriteChan chan<- dto.StrategyResponse) {
	for sub := range subscribers {
		go th.readFromSub(sub, wsWriteChan)
	}
}

func (th *TradingHandler) readFromSub(sub strategy.Subscription, wsWriteChan chan<- dto.StrategyResponse) {
	for msg := range sub.SubChan {
		strategyResponse := dto.StrategyResponse{
			StrategyMessage: dto.StrategyMessageResponse{
				Id:    msg.Id,
				Event: msg.Event,
			},
		}

		switch msg.Event {
		case "TASK_CREATION":
			mappedTask, err := mapper.MapTaskToReponseTask(msg.Task)
			if err != nil {
				logger.Error("failed to map task: ", err)
			}
			strategyResponse.StrategyMessage.Task = mappedTask

		case "STATUS_UPDATE":
			strategyResponse.StrategyMessage.State = msg.State

		case "MESSAGE_UPDATE":
			strategyResponse.StrategyMessage.Message = msg.Message
		}

		select {
		case wsWriteChan <- strategyResponse:
		default:
			logger.Error("Not able to push position response message to WS Write")
		}
	}
}

func (th *TradingHandler) writeToWs(ctx context.Context, ws <-chan dto.StrategyResponse, c *websocket.Conn) {
	for msg := range ws {
		err := wsjson.Write(ctx, c, msg)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (th *TradingHandler) handleError(ctx context.Context, err error, resp dto.StrategyResponse, c *websocket.Conn) {
	logger.Error(err.Error())
	resp.Error = err.Error()
	wsErr := wsjson.Write(ctx, c, resp)
	if wsErr != nil {
		logger.Error("error whilst writing: ", wsErr)
	}
}
