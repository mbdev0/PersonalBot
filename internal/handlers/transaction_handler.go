package handlers

import (
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions"
	"pump_fun/internal/webhook"
)

func HandleTransaction(signature string) {
	transactions.GetTransaction(signature)

	temp_coin := models.Coin{
		CoinData: models.MintData{
			Name:     "PumpFun",
			Symbol:   "PUMP",
			IPFS_URL: "https://ipfs.io/ipfs/QmZ4tZ",
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
