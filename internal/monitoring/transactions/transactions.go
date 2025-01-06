package transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mr-tron/base58"

	"pump_fun/internal/constants"
	"pump_fun/internal/launch"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions/decoder"

	"pump_fun/internal/logger"
)

func DecryptTransactionNotification(transaction models.TransactionNotification, coinStructChan chan models.Coin) *models.Coin {

	coin, transactionIsCreate := findCoinData(transaction) // we decrypt in here
	if !transactionIsCreate {
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

func findCoinData(transaction models.TransactionNotification) (coin models.Coin, transactionIsCreate bool) {

	coin = models.Coin{}
	transactionIsCreate = false

	for _, instruction := range transaction.Params.Result.Transaction.TransactionDetails.Message.Instructions {
		if len(instruction.Data) < 8 {
			continue
		}

		instructionData, err := base58.Decode(instruction.Data)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
			continue
		}

		discriminator := instructionData[:8]

		if bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:]) {
			err = decoder.DecodeCreateInstruction(&coin, instructionData)

			if err != nil {
				logger.Log(logger.LevelError, "Error decoding create instruction", logger.Error(err))
				continue
			}

			idlMap := launch.GetIdlMap()

			coin.CoinData.TokenAddr = instruction.Accounts[idlMap["create"].AccountMap["mint"]]
			coin.CoinData.CreatorAddr = instruction.Accounts[idlMap["create"].AccountMap["user"]]
			coin.CoinData.BondingCurveAddr = instruction.Accounts[idlMap["create"].AccountMap["bondingCurve"]]

			transactionIsCreate = true
		}

		if bytes.Equal(discriminator, constants.BuyInstructionDiscriminator[:]) {
			err := decoder.DecodeBuyInstruction(&coin, instructionData)

			if err != nil {
				logger.Log(logger.LevelError, "Error decoding buy instruction", logger.Error(err))
				continue
			}
		}
	}

	return coin, transactionIsCreate
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
