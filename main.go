// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/monitoring/transactions"
)

func main() {
	fmt.Println("Hello, World!")
	transaction := transactions.GetTransaction("5vQ5yPrGjE6LX5ZoLPCXK6YQtKymXLpeyVqBM6g2NUU86pvSSRVsTMG1FTyeTLPWErqSV8KAT2gD8bmEK6fzFVQg")
	fmt.Println(transaction)
	// monitoring.StartMonitor()
}
