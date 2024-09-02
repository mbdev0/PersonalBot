// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/config"
	"pump_fun/internal/monitoring"
	"pump_fun/internal/monitoring/transactions"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		return
	}

	config := config.GetConfig()
	fmt.Println("HTTP Node: ", config.HttpNode)
	fmt.Println("WS Node: ", config.WsNode)
	fmt.Println("Webhook: ", config.Webhook)

	monitoring.StartMonitor()

	transaction, err := transactions.GetTransaction("5vQ5yPrGjE6LX5ZoLPCXK6YQtKymXLpeyVqBM6g2NUU86pvSSRVsTMG1FTyeTLPWErqSV8KAT2gD8bmEK6fzFVQg")
	if err != nil {
		fmt.Println(err)
		return
	}
	decodedTransaction := transactions.DecodeCreateTransaction(transaction)
	fmt.Println(decodedTransaction)
}
