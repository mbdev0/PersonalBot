// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/buy"
	"pump_fun/internal/handlers"
	"pump_fun/internal/launch"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/models/tasks"

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

	tokenAddressPubkey, err := solana.PublicKeyFromBase58("6gFL19tSaRwtQCuQVEy6HykCbTHztbe1qxuodarypump")
	if err != nil {
		return
	}
	fmt.Println(tokenAddressPubkey.String())

	buyTask := tasks.BuyTask{
		Wallet:       privateKey,
		TokenAddress: tokenAddressPubkey,
		BuyAmount:    handlers.ConvertSolToLamport(0.001),
		Slippage:     0.2,
	}

	buy.SendBuyTransaction(&buyTask)
}
