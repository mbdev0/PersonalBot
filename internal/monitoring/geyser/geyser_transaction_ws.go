package geyser

import (
	"context"
	"pump_fun/internal/constants"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var (
	ws_url = config.GetConfig().WsNode
)

func Geyser_Stream_Transactions(transaction_chan chan<- models.TransactionNotification) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()
	defer func() {
		fmt.Println(("Exiting Geyser_Stream_Transactions"))
	}()
	defer close(transaction_chan)

	fmt.Println("Connecting to websocket at:", ws_url)
	ws, _, err := websocket.Dial(ctx, ws_url, nil)
	if err != nil {
		fmt.Println("Error connecting to websocket:", err)
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
			logger.Log(logger.LevelError, "Error reading from websocket", logger.Error(err))
			return err
		}

		transaction_chan <- out
	}
}
