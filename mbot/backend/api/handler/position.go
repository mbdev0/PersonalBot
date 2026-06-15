package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"personal_bot/backend/api/controller"
	"personal_bot/backend/api/dto"
	"personal_bot/backend/internal/services/subscription_hub/position"
	"personal_bot/backend/pkg/logger"
	"strconv"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
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
	mux.HandleFunc("GET /{id}", ph.getPositionById)
	mux.HandleFunc("GET /", ph.getPositions)
	mux.HandleFunc("GET /subscribe", ph.subscribe)
	mux.HandleFunc("GET /dashboard", ph.positionDashboard)
	mux.HandleFunc("DELETE /{id}", ph.delete)
}

func (ph *PositionHandler) positionDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := ph.controller.GetDashboard()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(dashboard)
	if err != nil {
		logger.Error(err)
	}
}

func (ph *PositionHandler) getPositionById(w http.ResponseWriter, r *http.Request) {
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

	position, err := ph.controller.GetBy(int64(convertedId))

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

func (ph *PositionHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	convId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID passed was not int", http.StatusBadRequest)
		return
	}

	err = ph.controller.Delete(r.Context(), int64(convId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *PositionHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	subscribers := make(chan position.Subscription, 1000)
	activeSubs := map[int64]position.Subscription{}

	defer close(subscribers)
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173", "localhost:9245"},
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
			ph.controller.Unsubscribe(strategyId, false)
		}
	}()

	wsWriteChan := make(chan dto.PositionResponse, 1000)
	defer close(wsWriteChan)

	//fan in loop
	go ph.fanIn(subscribers, wsWriteChan)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	//ws write loop
	go ph.writeToWs(ctx, wsWriteChan, c)

	for {
		var msg *dto.PositionSubscribe
		resp := dto.PositionResponse{}

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
			sub, err := ph.controller.Subscribe(msg.Id, false)
			if err != nil {
				ph.handleError(ctx, err, resp, c)
				continue
			}
			subscribers <- *sub
			activeSubs[sub.Sub_id] = *sub

		case dto.Unsubscribe:
			ph.controller.Unsubscribe(msg.Id, false)
			delete(activeSubs, msg.Id)
		}
	}

}

func (ph *PositionHandler) fanIn(subscribers <-chan position.Subscription, wsWriteChan chan<- dto.PositionResponse) {
	for sub := range subscribers {
		go ph.readFromSub(sub, wsWriteChan)
	}
}

func (ph *PositionHandler) readFromSub(sub position.Subscription, wsWriteChan chan<- dto.PositionResponse) {
	for msg := range sub.SubChan {
		posResp := dto.PositionResponse{
			PositionMessage: &msg,
		}
		select {
		case wsWriteChan <- posResp:
		default:
			logger.Error("Not able to push position response message to WS Write")
		}
	}
}

func (ph *PositionHandler) writeToWs(ctx context.Context, ws <-chan dto.PositionResponse, c *websocket.Conn) {
	for msg := range ws {
		err := wsjson.Write(ctx, c, msg)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (ph *PositionHandler) handleError(ctx context.Context, err error, resp dto.PositionResponse, c *websocket.Conn) {
	resp.Error = err.Error()
	wsErr := wsjson.Write(ctx, c, resp)
	if wsErr != nil {
		logger.Error(wsErr)
	}
}
