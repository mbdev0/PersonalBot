package geyser

import (
	"context"
	"pump_fun/internal/constants"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func Geyser_Stream_AccountInfo(ctx context.Context, address string, accountinfo_chan chan models.AccountSubscribeModel) error {

	ws, _, err := websocket.Dial(ctx, ws_url, nil)
	if err != nil {
		return err
	}

	ws.SetReadLimit(constants.WebSocketReadLimit)

	err = wsjson.Write(ctx, ws, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "accountSubscribe",
		"params": []interface{}{
			address,
			map[string]interface{}{
				"commitment": "processing",
				"encoding":   "base58",
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
		out := models.AccountSubscribeModel{}
		err = wsjson.Read(ctx, ws, &out)

		if err != nil {
			logger.Log(logger.LevelError, "Error reading from websocket", logger.Error(err))
			return err
		}

		accountinfo_chan <- out
	}
}
