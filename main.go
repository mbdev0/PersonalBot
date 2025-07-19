// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/handlers"
	"pump_fun/internal/launch"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/logger"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/transaction/buy"

	"github.com/gagliardetto/solana-go"
)

func main() {
	launch.LaunchOperations()
	// monitoring.StartAFKMonitor()

	privateKey, err := solana.PrivateKeyFromBase58(config.GetConfig().Wallet_Private_Key)
	if err != nil {
		logger.Error("Error creating private key from base58", err)
		return
	}

	tokenAddressPubkey, err := solana.PublicKeyFromBase58("Afk9Ms8AoUPbFzpGtLm4xx7m4UCAoGwhxuwcYNse4mjt")
	if err != nil {
		logger.Error("Error creating public key from base58", err)
		return
	}
	fmt.Println("Token Address PubKey:", tokenAddressPubkey.String())

	buyTask := tasks.BuyTask{
		Wallet:       privateKey,
		TokenAddress: tokenAddressPubkey,
		BuyAmount:    handlers.ConvertSolToLamport(0.001),
		Slippage:     0.2,
		BuyFee:       0.0001,
		ComputeUnits: 200000,
	}
	buyTask.InitDefaults()

	buy.SendBuyTransaction(&buyTask)

	fmt.Println("State: ", buyTask.State.TaskState.ToString())
	// sellTask := tasks.SellTask{
	// 	TokenAddress:     tokenAddressPubkey,
	// 	Wallet:           privateKey,
	// 	ComputeUnits:     100_000,
	// 	SellFee:          0.0001,
	// 	Slippage:         0.02,
	// 	PercentageToSell: 1,
	// }

	// sell.SendSellTransaction(&sellTask)
	// validationErrs := handlers.ValidateStruct(&buyTask)
	//
	//	if validationErrs != nil {
	//		logger.Error(validationErrs)
	//		return
	//	}
	//
	// //buy.SendBuyTransaction(&buyTask)
}
