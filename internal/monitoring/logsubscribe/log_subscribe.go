package logsubscribe

import (
	"context"
	"log/slog"
	"time"

	"pump_fun/internal/logger"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

const (
	ctxTimeout = 30 * time.Second
	ws_url     = "wss://mainnet.helius-rpc.com/"
)

func LogSubscribe() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	client, err := ws.Connect(ctx, ws_url)
	if err != nil {
		logger.Log(slog.LevelError, "Error connecting to websocket", slog.String("error", err.Error()))
		return err
	}

	// No returning errors as we need to keep a continuous connection to the websocket - we can log the errors instead
	pumpfunProgramId := solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	{
		sub, err := client.LogsSubscribeMentions(
			pumpfunProgramId,
			rpc.CommitmentConfirmed)

		if err != nil {
			logger.Log(slog.LevelError, "Error subscribing to logs", slog.String("error", err.Error()))
			return err
		}

		defer sub.Unsubscribe()

		for {
			msg, err := sub.Recv()
			if err != nil {
				logger.Log(slog.LevelError, "Error while streaming logs", slog.String("error", err.Error()))
			}

			if msg.Value.Err == nil {
				spew.Dump(msg.Value.Signature)
			}
		}
	}
}
