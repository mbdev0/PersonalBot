package stream

import (
	"context"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/monitoring/stream/response"
	"personal_bot/pkg/logger"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func GeyserStreamAccountInfo(ctx context.Context, address string, accountinfoChan chan response.AccountSubscribeModel) error {

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

	err = wsjson.Write(ctx, ws, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      420,
		"method":  "accountSubscribe",
		"params": []interface{}{
			address,
			map[string]interface{}{
				"commitment": "processed",
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
		out := response.AccountSubscribeModel{}
		err = wsjson.Read(ctx, ws, &out)

		if err != nil {
			logger.Error("Error reading from websocket", err)
			return err
		}

		accountinfoChan <- out
	}
}
