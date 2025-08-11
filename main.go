// temporary file for go build
package main

import (
	"net/http"
	"pump_fun/api/tasks/controller"
	"pump_fun/api/tasks/handler"
	"pump_fun/internal/launch"
	taskservice "pump_fun/internal/task_service"
	"pump_fun/internal/transaction"
	"pump_fun/pkg/logger"
)

func main() {
	launch.LaunchOperations()

	//TODO: move the init of api's into launch and return one mux back
	executor := transaction.TransactionExecutor{}
	taskService := taskservice.TaskService{Executor: &executor}
	taskService.NewTaskService()
	buyController := controller.TaskController{TaskService: &taskService}
	buyHandler := http.StripPrefix("/api/tasks", handler.NewTaskHandler(&buyController))

	mux := http.NewServeMux()
	mux.Handle("/api/tasks/", buyHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Information("Starting server on port 8080:")
	logger.Information("http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	// test()
}

func test() {
	// monitoring.StartAFKMonitor()

	// privateKey, err := solana.PrivateKeyFromBase58(config.GetConfig().Wallet_Private_Key)
	// if err != nil {
	// 	logger.Error("Error creating private key from base58", err)
	// 	return
	// }

	// tokenAddressPubkey, err := solana.PublicKeyFromBase58("Afk9Ms8AoUPbFzpGtLm4xx7m4UCAoGwhxuwcYNse4mjt")
	// if err != nil {
	// 	logger.Error("Error creating public key from base58", err)
	// 	return
	// }
	// fmt.Println("Token Address PubKey:", tokenAddressPubkey.String())

	// buyTask := tasks.BuyTask{
	// 	Wallet:       privateKey,
	// 	TokenAddress: tokenAddressPubkey,
	// 	BuyAmount:    handlers.ConvertSolToLamport(0.001),
	// 	Slippage:     0.2,
	// 	BuyFee:       0.0001,
	// 	ComputeUnits: 200000,
	// }
	// buyTask.InitDefaults()

	// buy.SendBuyTransaction(&buyTask)

	// fmt.Println("State: ", buyTask.State.TaskState.ToString())
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
