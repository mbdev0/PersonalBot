// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/launch"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/logger"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/sell"

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

	// tokenAddressPubkey, err := solana.PublicKeyFromBase58("6gFL19tSaRwtQCuQVEy6HykCbTHztbe1qxuodarypump")
	tokenAddressPubkey, err := solana.PublicKeyFromBase58("6gFL19tSaRwtQCuQVEy6HykCbTHztbe1qxuodarypump")
	if err != nil {
		logger.Error("Error creating public key from base58", err)
		return
	}
	fmt.Println("Token Address PubKey:", tokenAddressPubkey.String())

	// buyTask := tasks.BuyTask{
	// 	Wallet:       privateKey,
	// 	TokenAddress: tokenAddressPubkey,
	// 	BuyAmount:    handlers.ConvertSolToLamport(0.001),
	// 	Slippage:     0.2,
	// 	BuyFee:       0.0001,
	// 	ComputeUnits: 200000,
	// }

	// buy.SendBuyTransaction(&buyTask)

	sellTask := tasks.SellTask{
		PublicKey:    privateKey.PublicKey(),
		TokenAddress: tokenAddressPubkey,
		Wallet:       privateKey,
		ComputeUnits: 100_000,
		SellFee:      0.0001,
		Slippage:     0.02,
	}

	sell.SendSellTransaction(&sellTask)
}
