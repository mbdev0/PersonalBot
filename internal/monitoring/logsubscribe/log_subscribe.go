package logsubscribe

import (
	"context"
	"log/slog"

	"pump_fun/internal/logger"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func LogSubscribe() {
	client, err := ws.Connect(context.Background(), helius_ws)
	if err != nil {
		panic(err)
	}

	pumpfunProgramId := solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	{
		sub, err := client.LogsSubscribeMentions(
			pumpfunProgramId,
			rpc.CommitmentConfirmed)

		if err != nil {
			logger.Log(slog.LevelError, "Error subscribing to logs")
		}

		defer sub.Unsubscribe()

		for {
			msg, err := sub.Recv()
			if err != nil {
				logger.Log(slog.LevelError, "Error while streaming logs")
			}

			spew.Dump(msg)
		}
	}
}
