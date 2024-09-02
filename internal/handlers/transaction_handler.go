package handlers

import (
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions"
	"pump_fun/internal/webhook"
)

func HandleTransaction(signature string) {
	transaction, err := transactions.GetTransaction(signature)
	if err != nil {
		return
	}
	result := transactions.ParseTransaction(transaction)

	temp_coin := models.Coin{
		CoinData: models.MintData{
			Name:     result.Name,
			Symbol:   result.Symbol,
			IPFS_URL: result.IPFS_URL,
		},
		IPFSData: models.IPFS{
			TelegramURL: "https://t.me/pumpfun",
			TwitterURL:  "https://twitter.com/pumpfun",
			WebsiteURL:  "https://pump.fun",
			ImageURL:    "https://pump.fun/pumpfun.png",
		},
	}

	webhook.SendWebhook(&temp_coin)
}
