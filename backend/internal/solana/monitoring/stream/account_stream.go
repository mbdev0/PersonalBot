package stream

import (
	"context"
	"encoding/json"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/solana/monitoring/stream/response"
	"personal_bot/pkg/logger"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func GeyserStreamAccountInfo(ctx context.Context, address string, accountinfoChan chan response.AccountSubscribeModel, wsUrl string) error {

	ws, err := retry.DoWithData(
		func() (*websocket.Conn, error) {
			ws, _, err := websocket.Dial(ctx, wsUrl, nil)
			if err != nil {
				return nil, err
			}
			return ws, nil
		}, retry.Attempts(constants.Retries))

	if err != nil {
		return err
	}

	ws.SetReadLimit(constants.WebSocketReadLimit)

	err = wsjson.Write(ctx, ws, map[string]any{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "accountSubscribe",
		"params": []any{
			address,
			map[string]any{
				"commitment": "processed",
				"encoding":   "base64",
			},
		},
	})

	if err != nil {
		return err
	}

	var firstMessage any
	err = wsjson.Read(ctx, ws, &firstMessage)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Second)
	go func(ctx context.Context) {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ws.Ping(ctx)
			}
		}
	}(ctx)

	for {
		_, raw, err := ws.Read(ctx)
		if err != nil {
			logger.Error("Error reading from websocket", err)
			return err
		}

		out := response.AccountSubscribeModel{}
		if err := json.Unmarshal(raw, &out); err != nil {
			logger.Error("Error unmarshalling message", err, string(raw))
			continue
		}

		accountinfoChan <- out
	}
}

func NewGeyserStreamAccountInfo(ctx context.Context, address, wsUrl string, accountinfoChan chan<- []byte) error {

	ws, err := retry.DoWithData(
		func() (*websocket.Conn, error) {
			ws, _, err := websocket.Dial(ctx, wsUrl, nil)
			if err != nil {
				return nil, err
			}
			return ws, nil
		}, retry.Attempts(constants.Retries))

	if err != nil {
		return err
	}

	ws.SetReadLimit(constants.WebSocketReadLimit)

	err = wsjson.Write(ctx, ws, map[string]any{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "accountSubscribe",
		"params": []any{
			address,
			map[string]any{
				"commitment": "processed",
				"encoding":   "base64",
			},
		},
	})

	if err != nil {
		return err
	}

	var firstMessage any
	err = wsjson.Read(ctx, ws, &firstMessage)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Second)
	go func(ctx context.Context) {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ws.Ping(ctx)
			}
		}
	}(ctx)

	for {
		_, raw, err := ws.Read(ctx)
		if err != nil {
			logger.Error("Error reading from websocket", err)
			return err
		}

		accountinfoChan <- raw
	}
}
