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
