// temporary file for go build
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"pump_fun/api/controller"
	"pump_fun/api/handler"
	"pump_fun/app"
	"pump_fun/internal/services/state"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/internal/services/trading"
	"pump_fun/pkg/logger"
	"runtime"
	"time"
)

func main() {
	app.Launch()

	//TODO: move the init of api's into launch and return one mux back
	fsm := state.Machine{}
	subhub := subscriptionhub.Hub{}
	subhub.New()
	stateManger := state.Manager{}
	stateManger.New(&subhub)
	taskService := taskservice.TaskService{StateMachine: &fsm, StateManager: &stateManger, Hub: &subhub}
	taskService.NewTaskService()

	tradingStrategy := trading.Strategy{}
	tradingStrategy.NewTradingStrategy(&taskService)
	tradingService := trading.Service{}
	tradingService.NewTradingService(&tradingStrategy)

	tradingController := controller.StrategyController{}
	tradingController.New(&tradingService)
	tradingHandler := http.StripPrefix("/api/trading", handler.NewTradingHandler(&tradingController))

	buyController := controller.TaskController{TaskService: &taskService}
	buyHandler := http.StripPrefix("/api/tasks", handler.NewTaskHandler(&buyController))

	mux := http.NewServeMux()
	mux.Handle("/api/tasks/", buyHandler)
	mux.Handle("/api/trading/", tradingHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Information("Starting server on port 8080:")
	logger.Information("http://localhost:8080")

	// checkGoRoutines()
	// test()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}

func checkGoRoutines() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
		}
	}()

}

func test() {
	// monitoring.StartAFKMonitor()

	// privateKey, err := solana.PrivateKeyFromBase58(config.GetConfig().WalletPrivateKey)
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
