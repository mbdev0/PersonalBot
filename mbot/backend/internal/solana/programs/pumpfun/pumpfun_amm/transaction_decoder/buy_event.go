package transactiondecoder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	pumpfunamm "personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/constants"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/models"
	"strings"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func GetBuyEvent(tx *rpc.GetParsedTransactionResult) (models.BuyEvent, error) {
	if tx.Meta.Err != nil {
		return models.BuyEvent{}, fmt.Errorf("error from tx")
	}

	for _, log := range tx.Meta.LogMessages {
		if !strings.HasPrefix(log, "Program data: ") {
			continue
		}

		base64String := strings.TrimPrefix(log, "Program data: ")
		data, err := base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			continue
		}

		if len(data) < 8 {
			continue
		}

		if !bytes.Equal(data[:8], pumpfunamm.BuyEvent) {
			continue
		}

		var tradeEvent models.BuyEvent
		err = borsh.Deserialize(&tradeEvent, data[8:])
		if err != nil {
			continue
		}

		return tradeEvent, nil
	}

	return models.BuyEvent{}, fmt.Errorf("unable to find buy event in logs")
}
