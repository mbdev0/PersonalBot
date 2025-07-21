package geyser

import (
	"context"
	"pump_fun/internal/constants"
	"pump_fun/internal/models"
	"pump_fun/pkg/logger"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func Geyser_Stream_AccountInfo(ctx context.Context, address string, accountinfo_chan chan models.AccountSubscribeModel) error {

	ws, err := retry.DoWithData(
		func() (*websocket.Conn, error) {
			ws, _, err := websocket.Dial(ctx, ws_url, nil)
			if err != nil {
				return nil, err
			}
			return ws, nil
		}, retry.Attempts(constants.Retries))

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
				"encoding":   "base64",
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
			logger.Error("Error reading from websocket", err)
			return err
		}

		accountinfo_chan <- out
	}
}
