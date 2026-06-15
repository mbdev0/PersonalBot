package transactiondecoder

import (
	"encoding/base64"
	"fmt"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_native/models"
	"strings"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type InnerInstruction struct {
	ProgramId string `json:"program_id"`
}

func GetTradeEvent(tx *rpc.GetParsedTransactionResult) (models.TradeEvent, error) {
	if tx.Meta.Err != nil {
		return models.TradeEvent{}, fmt.Errorf("error from tx")
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

		var tradeEvent models.TradeEvent
		err = borsh.Deserialize(&tradeEvent, data[8:])
		if err != nil {
			continue
		}

		return tradeEvent, nil
	}

	return models.TradeEvent{}, fmt.Errorf("unable to find trade event in logs")
}
