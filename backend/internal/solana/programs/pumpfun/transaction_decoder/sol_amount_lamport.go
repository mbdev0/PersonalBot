package transactiondecoder

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func ExtractTotalSolSpent(tx *rpc.GetParsedTransactionResult, wallet solana.PublicKey) (float64, error) {
	transactionMessage := tx.Transaction.Message

	//extract sol amount
	walletIndex := -1

	for i, account := range transactionMessage.AccountKeys {
		if account.PublicKey == wallet {
			walletIndex = i
		}
	}

	if walletIndex == -1 {
		return 0, fmt.Errorf("could not find user's wallet in account keys")
	}

	solAmountLamport := tx.Meta.PreBalances[walletIndex] - tx.Meta.PostBalances[walletIndex]
	return float64(solAmountLamport), nil
}
