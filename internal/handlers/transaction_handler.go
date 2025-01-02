package handlers

import (
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions"
	"pump_fun/internal/monitoring/transactions/decoder"
)

func HandleTransactionNotification(decoder *decoder.Decoder, transaction models.TransactionNotification, coinStructChan chan models.Coin) {
	coin := transactions.DecryptTransactionNotification(decoder, transaction, coinStructChan)
	coinStructChan <- *coin
}
