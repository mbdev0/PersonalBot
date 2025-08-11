package handlers

import (
	wallet "pump_fun/internal/solana/wallet"

	"github.com/gagliardetto/solana-go"
)

func SignTx(tx *solana.Transaction, privateKey solana.PrivateKey) {
	wallet.SignTx(tx, privateKey)
}
