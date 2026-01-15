package decoder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mr-tron/base58"

	"personal_bot/app/pumpfun_idl"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/monitoring/models"
	"personal_bot/internal/monitoring/stream/response"

	"personal_bot/pkg/logger"
)

func DecryptTransactionNotificationForCoin(transaction response.TransactionNotification) *models.Coin {

	coin, err := getCreatedCoinWithBuyData(transaction)

	if err != nil {
		return nil
	}

	// we always get the IPFS data even when it's slow to get a response
	// TODO: ideally we should return basic information to the monitor first - access ipfs url if needed
	ipfsData, err := getIPFSData(coin.CoinData.IpfsUrl)
	if err != nil {
		logger.Error("Error getting IPFS data - IpfsUrl: ", coin.CoinData.IpfsUrl, " ", err)
		return nil
	}

	coin.CoinData.Signature = transaction.Params.Result.Signature
	coin.IPFSData = *ipfsData

	return &coin
}

func getCreatedCoinWithBuyData(transaction response.TransactionNotification) (models.Coin, error) {
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
			logger.Error("Error decoding instruction data", err)
			continue
		}

		discriminator := instructionData[:8]
		isCreateInstruction := bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:]) || bytes.Equal(discriminator, constants.CreateV2InstructionDiscriminator[:])
		isBuyInstruction := bytes.Equal(discriminator, constants.BuyInstructionDiscriminator[:])

		if isCreateInstruction {
			coin, err = createCoinFromInstruction(instruction, instructionData)
			if err != nil {
				logger.Error("Error creating coin from instruction", err)
				continue
			}
			createTransactionFound = true
		} else if isBuyInstruction {
			devHoldingAmount, err = ExtractBuyAmountFromBuyInstruction(instructionData)
			if err != nil {
				logger.Error("Error extracting buy amount from buy instruction", err)
				continue
			}
		}
	}

	if !createTransactionFound {
		return coin, errors.New("create instruction is not found")
	}

	coin.CoinData.DevHoldingAmount = devHoldingAmount
	return coin, nil
}

func createCoinFromInstruction(instruction response.Instruction, instructionData []byte) (models.Coin, error) {
	coin := models.Coin{}

	decodedInstruction, err := decodeCreateInstruction(instructionData)
	if err != nil {
		logger.Error("Error decoding create instruction ", err)
		return coin, err
	}
	updateCoinFromDecodedInstruction(&coin, decodedInstruction)
	assignCoinAddresses(&coin, instruction)

	return coin, nil
}

func assignCoinAddresses(coin *models.Coin, instruction response.Instruction) {
	createAccountIDL := pumpfun_idl.GetIdlMap()["create"].AccountMap

	coin.CoinData.TokenAddr = instruction.Accounts[createAccountIDL["mint"]]
	coin.CoinData.CreatorAddr = instruction.Accounts[createAccountIDL["user"]]
	coin.CoinData.BondingCurveAddr = instruction.Accounts[createAccountIDL["bondingCurve"]]
}

func getIPFSData(ipfsURL string) (*models.IPFS, error) {
	resp, err := http.Get(ipfsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ipfsData models.IPFS
	if err := json.NewDecoder(resp.Body).Decode(&ipfsData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &ipfsData, nil
}
