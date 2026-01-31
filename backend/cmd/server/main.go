// temporary file for go build
package main

import (
	"context"
	"net/http"
	"os/signal"
	"personal_bot/api/controller"
	"personal_bot/api/handler"
	"personal_bot/app"
	"personal_bot/infrastructure/persistence"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/services/position"
	"personal_bot/internal/services/state"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	positionhub "personal_bot/internal/services/subscription_hub/position"
	"personal_bot/internal/services/subscription_hub/strategy"
	"personal_bot/internal/services/wallet"
	"syscall"
	"time"

	taskservice "personal_bot/internal/services/task_service"
	"personal_bot/internal/services/trading"
	"personal_bot/internal/services/transaction"
	"personal_bot/pkg/logger"
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

	strategySubHub := strategy.NewSubscriptionHub()

	executor := transaction.Executor{}
	executor.New(&taskSubhub, &positionService)

	stateManger := state.Manager{}
	stateManger.New(&taskSubhub, &executor)

	taskRepo := repository.NewTaskRepository(db)
	taskService := taskservice.NewTaskService(&fsm, &stateManger, taskRepo, &taskSubhub)

	tradingStrategy := trading.Strategy{}
	tradingStrategy.NewTradingStrategy(taskService, &posSubhub, &positionService, &strategySubHub)

	tradingTaskRepo := repository.NewTradingRepository(db)
	tradingService := trading.Service{}
	tradingService.NewTradingService(&tradingStrategy, &strategySubHub, tradingTaskRepo, taskService)
	err = tradingService.LoadFromDB(context.Background())
	if err != nil {
		logger.Error(err)
	}

	err = taskService.LoadFromDB(context.Background())
	if err != nil {
		logger.Error("error whilst loading from db: ", err)
	}

	walletRepo := repository.NewWalletRepository(db)
	walletService := wallet.NewWalletService(walletRepo)
	walletController := controller.NewWalletController(walletService)
	walletHandler := http.StripPrefix("/api/wallet", handler.NewWalletHandler(walletController))

	tradingController := controller.StrategyController{}
	tradingController.New(&tradingService, walletService)
	tradingHandler := http.StripPrefix("/api/trading", handler.NewTradingHandler(&tradingController))

	buyController := controller.TaskController{TaskService: taskService, WalletService: walletService}
	buyHandler := http.StripPrefix("/api/tasks", handler.NewTaskHandler(&buyController))

	positionController := controller.PositionController{PositionService: &positionService}
	positionHandler := http.StripPrefix("/api/position", handler.NewPositionHandler(&positionController))

	dashboardHandler := http.StripPrefix("/api/dashboard", handler.NewDashboardHandler(&tradingController, &buyController))

	mux := http.NewServeMux()
	mux.Handle("/api/tasks/", buyHandler)
	mux.Handle("/api/trading/", tradingHandler)
	mux.Handle("/api/position/", positionHandler)
	mux.Handle("/api/wallet/", walletHandler)
	mux.Handle("/api/dashboard/", dashboardHandler)

	server := &http.Server{
		Addr:    ":9090",
		Handler: CORS(mux.ServeHTTP),
	}

	// checkGoRoutines()
	// test()

	go func() {
		logger.Information("Starting server on port 9090:")
		logger.Information("http://localhost:9090")
		if err := server.ListenAndServe(); err != nil {
			logger.Information("stopped listening to localhost:9090")
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdown.Done()

	logger.Information("shutting down application")

	logger.Information("pushing tasks to db")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = tradingService.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}

	// do the tasks here
	err = taskService.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("error whilst shutting down: ", err)
	}

	err = db.Close()
	if err != nil {
		logger.Error(err)
	}

	logger.Information("shutdown complete")

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
