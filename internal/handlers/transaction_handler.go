package handlers

import (
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions"
)

func HandleTransactionNotification(transaction models.TransactionNotification, coinStructChan chan models.Coin) {
	coin := transactions.DecryptTransactionNotification(transaction, coinStructChan)

	if coin != nil {
		coin = HandleCoinFiltering(coin)
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
