package handlers

import (
	"fmt"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions"
)

func HandleTransaction(signature string) {
	result, err := transactions.GetTransaction(signature)
	if err != nil {
		return
	}
	ipfsData, err := transactions.GetIPFSData(result.IPFS_URL)

	if err != nil {
		fmt.Println("Error getting IPFS data", err)
		return
	}

	tempCoin := models.Coin{
		CoinData: models.MintData{
			Name:     result.Name,
			Symbol:   result.Symbol,
			IPFS_URL: result.IPFS_URL,
		},
		IPFSData: *ipfsData,
	}

	fmt.Print("temp_coin\n", tempCoin, "\n\n")
	// webhook.SendWebhook(&temp_coin)
}
