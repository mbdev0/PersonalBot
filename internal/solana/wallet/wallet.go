package wallets

import (
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

func SignTx(tx *solana.Transaction, privateKey solana.PrivateKey) {

	_, err := tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privateKey.PublicKey().Equals(key) {
				return &privateKey
			}
			return nil
		},
	)

	if err != nil {
		logger.Error("Error signing transaction", err)
	}
}
