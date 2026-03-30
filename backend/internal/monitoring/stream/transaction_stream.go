package stream

import (
	"context"
	"fmt"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/monitoring/stream/response"
	"personal_bot/pkg/logger"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func StartGeyserTransactionStream(transactionChan chan<- response.TransactionNotification, ctx context.Context, wsUrl string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := geyserStreamTransactions(transactionChan, ctx, wsUrl)
			if err != nil {
				return err
			}
		}
	}
}

func geyserStreamTransactions(transactionChan chan<- response.TransactionNotification, ctx context.Context, wsUrl string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), connection_timeout)
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

	err = wsjson.Write(ctx, ws, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "transactionSubscribe",
		"params": []interface{}{
			map[string]interface{}{
				"failed": false,
				"accountInclude": []interface{}{
					constants.Program,
				},
			},
			map[string]interface{}{
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

	var firstMessage interface{}
	err = wsjson.Read(ctx, ws, &firstMessage)
	if err != nil {
		return err
	}
	logger.Information("connected to ws")

	for {
		out := response.TransactionNotification{}
		err = wsjson.Read(ctx, ws, &out)

		if err != nil {
			logger.Error("Error reading from websocket: ", err)
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case transactionChan <- out:
		}

	}
}
