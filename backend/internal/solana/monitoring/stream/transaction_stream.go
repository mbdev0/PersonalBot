package stream

import (
	"context"
	"fmt"
	"personal_bot/internal/core/constants"
	"personal_bot/pkg/logger"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func NewStartGeyserTransactionStream(ctx context.Context, program, wsUrl string, transactionChan chan<- []byte) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := newGeyserStreamTransactions(ctx, transactionChan, program, wsUrl)
			if err != nil {
				return err
			}
		}
	}
}

func newGeyserStreamTransactions(ctx context.Context, transactionChan chan<- []byte, program, wsUrl string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ws, err := retry.DoWithData(
		func() (*websocket.Conn, error) {
			fmt.Println("Connecting to websocket...")
			ws, _, err := websocket.Dial(ctx, wsUrl, nil)
			if err != nil {
				return nil, err
			}
			return ws, nil
		}, retry.Attempts(constants.Retries))

	if err != nil {
		return err
	}

	defer func(ws *websocket.Conn, code websocket.StatusCode, reason string) {
		err := ws.Close(code, reason)
		if err != nil {
			logger.Error(err.Error())
		}
	}(ws, websocket.StatusNormalClosure, "websocket closed")

	ws.SetReadLimit(constants.WebSocketReadLimit)

	err = wsjson.Write(ctx, ws, map[string]any{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "transactionSubscribe",
		"params": []any{
			map[string]any{
				"failed": false,
				"accountInclude": []any{
					program,
				},
			},
			map[string]any{
				"commitment":                     "confirmed",
				"transactionDetails":             "full",
				"encoding":                       "jsonParsed",
				"maxSupportedTransactionVersion": 0,
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
	logger.Information("connected to ws")

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
		_, data, err := ws.Read(ctx)

		if err != nil {
			logger.Error("Error reading from websocket: ", err)
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case transactionChan <- data:
		}

	}
}
