package client

import (
	"context"
	"fmt"
	"pump_fun/pkg/logger"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	contextTimeout   = 30 * time.Second
	pollInterval     = 2 * time.Second
	maxConfirmations = 31
)

func ConfirmTransaction(sig solana.Signature, ctx context.Context) (IsSuccess bool, err error) {
	client := GetClient()
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	for {
		txResp, err := client.GetSignatureStatuses(ctx, true, sig)
		if err != nil {
			return false, err
		}

		if len(txResp.Value) == 0 || txResp.Value[0] == nil {
			logger.Information("Transaction status unavailable. Retrying...")
			time.Sleep(pollInterval)
			continue
		}

		status := txResp.Value[0]
		if status.Err != nil {
			return false, fmt.Errorf("transaction failed: %v", status.Err)
		}

		switch status.ConfirmationStatus {
		case rpc.ConfirmationStatusFinalized:
			return true, nil
		case rpc.ConfirmationStatusConfirmed:
			logger.Information(fmt.Sprintf("Transaction confirmed: %d/%d confirmations",
				*status.Confirmations, maxConfirmations))

			expired, err := IsBlockhashExpired(status.Slot, ctx)
			if err != nil {
				return false, fmt.Errorf("blockhash expiration check failed: %w", err)
			}
			if expired {
				return false, fmt.Errorf("blockhash expired, re-send the tx")
			}

			time.Sleep(pollInterval)
		}
	}
}

type ConfirmMessage struct {
	Message string
	Err     string
}

func ConfirmTransactionWithStream(sig solana.Signature, ctx context.Context, stream chan ConfirmMessage) {
	client := GetClient()
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	for {
		msg := ConfirmMessage{}
		txResp, err := client.GetSignatureStatuses(ctx, true, sig)
		if err != nil {
			msg.Err = err.Error()
			stream <- msg
			return
		}

		if len(txResp.Value) == 0 || txResp.Value[0] == nil {
			logger.Information("Transaction status unavailable. Retrying...")
			msg.Message = "Transaction status unavailable. Retrying..."
			stream <- msg
			time.Sleep(pollInterval)
			continue
		}

		status := txResp.Value[0]
		if status.Err != nil {
			msg.Err = fmt.Sprintf("transaction failed: %v", status.Err)
			stream <- msg
			return
		}

		switch status.ConfirmationStatus {
		case rpc.ConfirmationStatusFinalized:
			msg.Message = fmt.Sprintf("Succesfully confirmed tx: %s", sig)
			stream <- msg
			return
		case rpc.ConfirmationStatusConfirmed:
			logger.Information(fmt.Sprintf("Transaction confirmed: %d/%d confirmations",
				*status.Confirmations, maxConfirmations))

			msg.Message = fmt.Sprintf("Transaction confirmed: %d/%d confirmations",
				*status.Confirmations, maxConfirmations)
			stream <- msg

			expired, err := IsBlockhashExpired(status.Slot, ctx)
			if err != nil {
				msg.Err = fmt.Sprintf("blockhash expiration check failed: %w", err)
				stream <- msg
				return
			}
			if expired {
				msg.Err = "blockhash expired, re-send the tx"
				stream <- msg
				return
			}

			time.Sleep(pollInterval)
		}
	}
}
