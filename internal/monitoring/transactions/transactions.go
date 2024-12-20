package transactions

import (
	"bytes"

	"github.com/mr-tron/base58"

	"pump_fun/internal/constants"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions/decoder"

	"pump_fun/internal/logger"
)

func DecryptTransactionNotification(decoder *decoder.Decoder, transaction models.TransactionNotification, coinStructChan chan models.Coin) {

	// accounts := transaction.Params.Result.Transaction.TransactionDetails.Message.AccountKeys
	compiled_instruction := findInstruction(transaction)

	if len(compiled_instruction.Data) < 8 {
		return
	}

	//decode instruction from base58
	instruction, err := base58.Decode(compiled_instruction.Data)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
		return
	}

	discriminator := instruction[:8]

	if bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:]) {
		decodedInstruction, err := decoder.Decode(instruction)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
			return
		}
		coinStructChan <- mapToStruct(decodedInstruction)
	}

}

func findInstruction(transaction models.TransactionNotification) models.Instruction {
	for _, instruction := range transaction.Params.Result.Transaction.TransactionDetails.Message.Instructions {
		if len(instruction.Data) < 8 {
			continue
		}

		instructionData, err := base58.Decode(instruction.Data)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
			return models.Instruction{}
		}

		discriminator := instructionData[:8]

		if bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:]) {
			return instruction
		}
	}

	return models.Instruction{}
}

func mapToStruct(decodedInstruction models.DecodedInstruction) models.Coin {
	return models.Coin{
		CoinData: models.MintData{
			Name:             decodedInstruction.Name,
			Symbol:           decodedInstruction.Symbol,
			IPFS_URL:         decodedInstruction.IPFS_URL,
			TokenAddr:        "",
			CreatorAddr:      "",
			DevHoldingAmount: 0,
		},
		IPFSData: models.IPFS{
			TelegramURL: "https://t.me/pumpfun",
			TwitterURL:  "https://twitter.com",
			WebsiteURL:  "https://pump.fun",
			ImageURL:    "https://pump.fun/pumpfun.png",
		},
	}
}
