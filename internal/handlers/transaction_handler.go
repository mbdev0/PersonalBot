package handlers

import (
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/decoder"
)

func HandleTransactionNotification(transaction models.TransactionNotification, coinStructChan chan<- models.Coin) {
	coin := decoder.DecryptTransactionNotificationForCoin(transaction)

	if coin != nil {
		coin = HandleCoinFiltering(coin)
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
