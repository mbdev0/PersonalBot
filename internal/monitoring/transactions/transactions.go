package transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mr-tron/base58"

	"pump_fun/internal/constants"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions/decoder"

	"pump_fun/internal/logger"
)

func DecryptTransactionNotification(decoder *decoder.Decoder, transaction models.TransactionNotification, coinStructChan chan models.Coin) *models.Coin {

	// accounts := transaction.Params.Result.Transaction.TransactionDetails.Message.AccountKeys
	compiled_instruction := findInstruction(transaction)
	if len(compiled_instruction.Data) < 8 {
		return nil
	}

	//decode instruction from base58
	instruction, err := base58.Decode(compiled_instruction.Data)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
		return nil
	}

	discriminator := instruction[:8]

	if bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:]) {
		decodedInstruction, err := decoder.Decode(instruction)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
			return nil
		}

		ipfsData, err := GetIPFSData(decodedInstruction.IPFS_URL)
		if err != nil {
			logger.Log(logger.LevelError, "Error getting IPFS data", logger.Error(err))
			return nil
		}

		coin := mapToStruct(decodedInstruction, *ipfsData)
		return &coin
	}
	return nil
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

func mapToStruct(decodedInstruction models.DecodedInstruction, ipfs models.IPFS) models.Coin {
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
			TelegramURL: ipfs.TelegramURL,
			TwitterURL:  ipfs.TwitterURL,
			WebsiteURL:  ipfs.WebsiteURL,
			ImageURL:    ipfs.ImageURL,
		},
	}
}

func GetIPFSData(ipfsURL string) (*models.IPFS, error) {
	resp, err := http.Get(ipfsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ipfsData models.IPFS
	if err := json.NewDecoder(resp.Body).Decode(&ipfsData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &ipfsData, nil
}
