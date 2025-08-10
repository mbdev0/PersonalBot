package geyser

import (
	"context"
	"fmt"
	"pump_fun/internal/constants"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/models"
	"pump_fun/pkg/logger"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var (
	ws_url             = config.GetConfig().WsNode
	connection_timeout = time.Second * 10
)

func Geyser_Stream_Transactions(transaction_chan chan<- models.TransactionNotification) error {
	ctx, cancel := context.WithTimeout(context.Background(), connection_timeout)
	defer cancel()

	ws, err := retry.DoWithData(
		func() (*websocket.Conn, error) {
			fmt.Println("Connecting to websocket...")
			ws, _, err := websocket.Dial(ctx, ws_url, nil)
			if err != nil {
				return nil, err
			}
			return ws, nil
		}, retry.Attempts(constants.Retries))

	if err != nil {
		return err
	}

	defer ws.Close(websocket.StatusNormalClosure, "websocket closed")
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

	for {
		out := models.TransactionNotification{}
		err = wsjson.Read(ctx, ws, &out)

		if err != nil {
			logger.Error("Error reading from websocket", err)
			return err
		}

		transaction_chan <- out
	}
}
