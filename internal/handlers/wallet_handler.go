package handlers

import (
	"pump_fun/internal/wallets"

	"github.com/gagliardetto/solana-go"
)

func SignTx(tx *solana.Transaction, privateKey solana.PrivateKey) {
	wallets.SignTx(tx, privateKey)
}
