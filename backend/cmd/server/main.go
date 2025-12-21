// temporary file for go build
package main

import (
	"net/http"
	"pump_fun/api/controller"
	"pump_fun/api/handler"
	"pump_fun/app"
	"pump_fun/infrastructure/persistence"
	"pump_fun/infrastructure/persistence/repository"
	"pump_fun/internal/services/position"
	"pump_fun/internal/services/state"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	positionhub "pump_fun/internal/services/subscription_hub/position"
	"pump_fun/internal/services/wallet"

	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/internal/services/trading"
	"pump_fun/internal/services/transaction"
	"pump_fun/pkg/logger"
)

func main() {
	db, err := persistence.NewConnection()
	if err != nil {
		panic(err)
	}

	app.Launch()

	//TODO: move the init of api's into launch and return one mux back
	fsm := state.Machine{}
	taskSubhub := subscriptionhub.Hub{}
	taskSubhub.New()

	posSubhub := positionhub.NewSubscriptionHub()
	positionService := position.NewPositionService(&posSubhub)

	executor := transaction.Executor{}
	executor.New(&taskSubhub, &positionService)

	stateManger := state.Manager{}
	stateManger.New(&taskSubhub, &executor)

	taskService := taskservice.TaskService{StateMachine: &fsm, StateManager: &stateManger, Hub: &taskSubhub}
	taskService.NewTaskService()

	walletRepo := repository.NewWalletRepository(db)
	walletService := wallet.NewWalletService(walletRepo)
	walletController := controller.NewWalletController(walletService)
	walletHandler := http.StripPrefix("/api/wallet", handler.NewWalletHandler(walletController))

	tradingStrategy := trading.Strategy{}
	tradingStrategy.NewTradingStrategy(&taskService, &posSubhub, &positionService)

	tradingService := trading.Service{}
	tradingService.NewTradingService(&tradingStrategy)

	tradingController := controller.StrategyController{}
	tradingController.New(&tradingService, walletService)
	tradingHandler := http.StripPrefix("/api/trading", handler.NewTradingHandler(&tradingController))

	buyController := controller.TaskController{TaskService: &taskService, WalletService: walletService}
	buyHandler := http.StripPrefix("/api/tasks", handler.NewTaskHandler(&buyController))

	positionController := controller.PositionController{PositionService: &positionService}
	positionHandler := http.StripPrefix("/api/position", handler.NewPositionHandler(&positionController))

	mux := http.NewServeMux()
	mux.Handle("/api/tasks/", buyHandler)
	mux.Handle("/api/trading/", tradingHandler)
	mux.Handle("/api/position/", positionHandler)
	mux.Handle("/api/wallet/", walletHandler)

	server := &http.Server{
		Addr:    ":9090",
		Handler: CORS(mux.ServeHTTP),
	}

	logger.Information("Starting server on port 9090:")
	logger.Information("http://localhost:9090")

	// checkGoRoutines()
	// test()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// func checkGoRoutines() {
// 	go func() {
// 		err := http.ListenAndServe("localhost:6060", nil)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			time.Sleep(3 * time.Second)
// 			fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
// 		}
// 	}()

// }

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
