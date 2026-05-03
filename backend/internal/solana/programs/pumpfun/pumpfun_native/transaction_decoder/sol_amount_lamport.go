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
			break
		}
	}

	if walletIndex == -1 {
		return 0, fmt.Errorf("could not find user's wallet in account keys")
	}

	solAmountLamport := int64(tx.Meta.PostBalances[walletIndex]) - int64(tx.Meta.PreBalances[walletIndex])
	return float64(solAmountLamport), nil
}
