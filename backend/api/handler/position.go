package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pump_fun/api/controller"
	"pump_fun/api/dto"
	"pump_fun/internal/services/subscription_hub/position"
	"pump_fun/pkg/logger"
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
	mux.HandleFunc("GET /position/{id}", ph.getPositionById)
	mux.HandleFunc("GET /positions", ph.getPositions)
	mux.HandleFunc("/subscribe", ph.subscribe)
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

func (ph *PositionHandler) subscribe(w http.ResponseWriter, r *http.Request) {
	subscribers := make(chan position.Subscription, 1000)
	defer close(subscribers)

	for k, v := range r.Header {
		log.Printf("Header: %s = %v", k, v)
	}

	//upgrade to ws
	c, err := websocket.Accept(w, r, nil)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error whilst trying to transition to WS", http.StatusInternalServerError)
		return
	}

	defer func(c *websocket.Conn) {
		err := c.CloseNow()
		if err != nil {
			logger.Error(err.Error())
		}
	}(c)

	wsWriteChan := make(chan dto.PositionResponse, 1000)
	defer close(wsWriteChan)

	//fan in loop
	go ph.fanIn(subscribers, wsWriteChan)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	//ws write loop
	go ph.writeToWs(wsWriteChan, c, ctx)

	for {
		var msg *dto.PositionSubscribe
		resp := dto.PositionResponse{}

		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			if websocket.CloseStatus(err) != -1 || ctx.Err() != nil {
				logger.Information("WebSocket connection closed or context cancelled")
				return // Clean exit
			}
			ph.handleError(err, resp, c, ctx)
			return
		}
		fmt.Println(msg)

		switch msg.Type {
		case dto.Subscribe:
			sub, err := ph.controller.Subscribe(msg.Id, false)
			if err != nil {
				ph.handleError(err, resp, c, ctx)
				continue
			}
			subscribers <- *sub

		case dto.Unsubscribe:
			err := ph.controller.Unsubscribe(msg.Id, false)
			if err != nil {
				ph.handleError(err, resp, c, ctx)
				continue
			}
		}
	}

	//go func() -> for every sub in ph.subs -> go func() -> write to wsWriteChan
	//ws write loop (in go func)
	//	for every msg -> write to ws
	//ws read loop (blocking)
	//	if sub -> call ph.posservice.sub, add to ph.subs => chan of subs
	//	if unsub ->  call ph.posservice.unsub, remove from ph.subs + call cancelCtx to kill go routine
}

func (ph *PositionHandler) fanIn(subscribers <-chan position.Subscription, wsWriteChan chan<- dto.PositionResponse) {
	for sub := range subscribers {
		fmt.Println("subscribers!")
		go ph.readFromSub(sub, wsWriteChan)
	}
}

func (ph *PositionHandler) readFromSub(sub position.Subscription, wsWriteChan chan<- dto.PositionResponse) {
	for msg := range sub.SubChan {
		fmt.Println("message from sub!")
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

func (ph *PositionHandler) writeToWs(ws <-chan dto.PositionResponse, c *websocket.Conn, ctx context.Context) {
	for msg := range ws {
		err := wsjson.Write(ctx, c, msg)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (ph *PositionHandler) handleError(err error, resp dto.PositionResponse, c *websocket.Conn, ctx context.Context) {
	resp.Error = err.Error()
	wsErr := wsjson.Write(ctx, c, resp)
	if wsErr != nil {
		logger.Error(wsErr)
	}
}
