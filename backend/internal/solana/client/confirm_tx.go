package client

import (
	"context"
	"fmt"
	"personal_bot/pkg/logger"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	contextTimeout   = 30 * time.Second
	pollInterval     = 2 * time.Second
	maxConfirmations = 32
)

type ConfirmMessage struct {
	Message string
	Err     string
}

func ConfirmTransactionWithStream(ctx context.Context, rpcClient *rpc.Client, sig solana.Signature, stream chan ConfirmMessage) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			stream <- ConfirmMessage{
				Err: "transaction confirmation timeout",
			}
			return
		default:
		}

		txResp, err := rpcClient.GetSignatureStatuses(ctx, true, sig)
		if err != nil {
			stream <- ConfirmMessage{
				Err: fmt.Sprintf("failed to get signature status: %v", err),
			}
			return
		}

		if len(txResp.Value) == 0 || txResp.Value[0] == nil {
			logger.Information("transaction status unavailable. Retrying...")
			stream <- ConfirmMessage{
				Message: "transaction status unavailable, retrying...",
			}
			time.Sleep(pollInterval)
			continue
		}

		status := txResp.Value[0]
		if status.Err != nil {
			logger.Error(status.Err)
			logger.Information(status)
			stream <- ConfirmMessage{
				Err: fmt.Sprintf("transaction failed: %v", status.Err),
			}
			return
		}

		switch status.ConfirmationStatus {
		case rpc.ConfirmationStatusFinalized:
			logger.Information(ConfirmMessage{
				Message: fmt.Sprintf("successfully confirmed transaction: %s", sig.String()),
			})

			stream <- ConfirmMessage{
				Message: fmt.Sprintf("successfully confirmed transaction: %s", sig.String()),
			}
			return
		case rpc.ConfirmationStatusConfirmed:
			confirmations := uint64(0)
			if status.Confirmations != nil {
				confirmations = *status.Confirmations
			}

			logger.Information(fmt.Sprintf("transaction confirmed: %d/%d confirmations",
				confirmations, maxConfirmations))

			stream <- ConfirmMessage{
				Message: fmt.Sprintf("transaction confirmed: %d/%d confirmations",
					confirmations, maxConfirmations),
			}

			expired, err := IsBlockhashExpired(ctx, status.Slot, rpcClient)
			if err != nil {
				stream <- ConfirmMessage{
					Err: fmt.Sprintf("blockhash expiration check failed: %v", err),
				}
				return
			}
			if expired {
				stream <- ConfirmMessage{
					Err: "blockhash expired, please re-send the transaction",
				}
				return
			}

			time.Sleep(pollInterval)

		case rpc.ConfirmationStatusProcessed:
			logger.Information("transaction processed, waiting for confirmation...")
			stream <- ConfirmMessage{
				Message: "transaction processed, waiting for confirmation...",
			}
			time.Sleep(pollInterval)

		default:
			logger.Information(fmt.Sprintf("unknown confirmation status: %v", status.ConfirmationStatus))
			stream <- ConfirmMessage{
				Message: fmt.Sprintf("unknown status: %v, retrying...", status.ConfirmationStatus),
			}
			time.Sleep(pollInterval)
		}
	}
}
