package geyser

import (
	"context"
	"fmt"
	"pump_fun/internal/constants"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/models"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var (
	ws_url = config.GetConfig().WsNode
)

func Geyser_Stream_Transactions(transaction_chan chan models.TransactionNotification) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws, _, err := websocket.Dial(ctx, ws_url, nil)
	if err != nil {
		return err
	}

	ws.SetReadLimit(constants.WebSocketReadLimit)

	err = wsjson.Write(ctx, ws, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "transactionSubscribe",
		"params": []interface{}{
			map[string]interface{}{
				"failed": false,
				"accountInclude": []interface{}{
					constants.ProgramID,
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
			fmt.Println("Error reading from websocket:", err)
			return err
		}

		transaction_chan <- out
	}
}
