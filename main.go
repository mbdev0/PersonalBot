// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/buy"
	"pump_fun/internal/handlers"
	"pump_fun/internal/launch"
	"pump_fun/internal/launch/config"

	"github.com/gagliardetto/solana-go"
)

func main() {
	launch.LaunchOperations()
	// monitoring.StartAFKMonitor()

	privateKey, err := solana.PrivateKeyFromBase58(config.GetConfig().Wallet_Private_Key)
	if err != nil {
		return
	}

	fmt.Println(privateKey.PublicKey().String())

	buy.SendBuyTransaction(privateKey.String(), "C2Fs7X9KpivZTPFr26ur8PR1iqVTNoSoZ7FuK2q1pump", handlers.ConvertSolToLamport(0.1), 0.2)
}
