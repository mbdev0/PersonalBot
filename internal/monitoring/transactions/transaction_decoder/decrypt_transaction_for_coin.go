package transaction_decoder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mr-tron/base58"

	"pump_fun/internal/constants"
	"pump_fun/internal/launch/pumpfun_idl"
	"pump_fun/internal/models"

	"pump_fun/internal/logger"
)

func DecryptTransactionNotificationForCoin(transaction models.TransactionNotification, coinStructChan chan models.Coin) *models.Coin {

	coin, err := getCreatedCoinWithBuyData(transaction)

	if err != nil {
		return nil
	}

	ipfsData, err := GetIPFSData(coin.CoinData.IPFS_URL)
	if err != nil {
		logger.Log(logger.LevelError, "Error getting IPFS data", logger.Error(err), logger.String("IPFS URL", coin.CoinData.IPFS_URL))
		return nil
	}

	coin.CoinData.Signature = transaction.Params.Result.Signature
	coin.IPFSData = *ipfsData

	return &coin
}

func getCreatedCoinWithBuyData(transaction models.TransactionNotification) (models.Coin, error) {
	var coin models.Coin
	var createTransactionFound bool
	var devHoldingAmount float64

	instructions := transaction.Params.Result.Transaction.TransactionDetails.Message.Instructions
	for _, instruction := range instructions {
		if len(instruction.Data) < 8 {
			continue
		}

		instructionData, err := base58.Decode(instruction.Data)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
			continue
		}

		discriminator := instructionData[:8]
		isCreateInstruction := bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:])
		isBuyInstruction := bytes.Equal(discriminator, constants.BuyInstructionDiscriminator[:])

		if isCreateInstruction {
			coin, err = createCoinFromInstruction(instruction, instructionData)
			if err != nil {
				logger.Log(logger.LevelError, "Error creating coin from instruction", logger.Error(err))
				continue
			}
			createTransactionFound = true
		} else if isBuyInstruction {
			devHoldingAmount, err = ExtractBuyAmountFromBuyInstruction(instructionData)
			if err != nil {
				logger.Log(logger.LevelError, "Error fetching buy amount from buy instruction", logger.Error(err))
				continue
			}
		}
	}

	if !createTransactionFound {
		return models.Coin{}, errors.New("create instruction is not found")
	}

	coin.CoinData.DevHoldingAmount = devHoldingAmount
	return coin, nil
}

func createCoinFromInstruction(instruction models.Instruction, instructionData []byte) (models.Coin, error) {
	coin := models.Coin{}

	decodedInstruction, err := DecodeCreateInstruction(instructionData)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding create instruction", logger.Error(err))
		return coin, err
	}
	UpdateCoinFromDecodedInstruction(&coin, decodedInstruction)
	assignCoinAddresses(&coin, instruction)

	return coin, nil
}

func assignCoinAddresses(coin *models.Coin, instruction models.Instruction) {
	idlMap := pumpfun_idl.GetIdlMap()
	createAccountIDL := idlMap["create"].AccountMap

	coin.CoinData.TokenAddr = instruction.Accounts[createAccountIDL["mint"]]
	coin.CoinData.CreatorAddr = instruction.Accounts[createAccountIDL["user"]]
	coin.CoinData.BondingCurveAddr = instruction.Accounts[createAccountIDL["bondingCurve"]]
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
