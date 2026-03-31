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
	cryptostates "personal_bot/internal/core/tasks/crypto_states"
	"personal_bot/internal/services/notifier"
	"personal_bot/internal/services/position"
	rpcgroups "personal_bot/internal/services/rpc_groups"
	"personal_bot/internal/services/settings"
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

	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	db, err := persistence.NewConnection()
	if err != nil {
		panic(err)
	}

	//TODO: move the init of api's into launch and return one mux back

	settingsRepo := repository.NewSettingRepository(db)
	settingsService := settings.NewSettingsService(settingsRepo)
	settings, err := settingsService.GetSettings(context.Background())
	if err != nil {
		logger.Error(err)
	}

	app.Launch(rpc.New(settings.PositionNodes.HTTPNode))

	discordNotifier := notifier.NewDiscordNotifier(settings)
	settingsController := controller.NewSettingsController(settingsService, discordNotifier)
	settingsHandler := http.StripPrefix("/api/settings", handler.NewSettingsHandler(settingsController))

	taskSubhub := subscriptionhub.NewTaskSubscriptionHub()
	posSubhub := positionhub.NewSubscriptionHub(settingsService)
	positionService := position.NewPositionService(posSubhub)

	rpcGroupRepo := repository.NewRPCGroupRepository(db)
	rpcGroupService := rpcgroups.NewRPCGroupService(rpcGroupRepo)
	rpcGroupController := controller.NewRPCGroupController(rpcGroupService)
	rpcGroupHandler := http.StripPrefix("/api/rpc_groups", handler.NewRPCGroupHandler(rpcGroupController))

	strategySubHub := strategy.NewSubscriptionHub()

	deps := cryptostates.Dependencies{
		Publisher:       taskSubhub,
		PositionService: positionService,
		Notifier:        discordNotifier,
	}

	fsmSteps := cryptostates.Build(&deps)
	executor := transaction.NewExecutor(taskSubhub, positionService, fsmSteps)

	taskRepo := repository.NewTaskRepository(db)
	taskManager := taskservice.NewTaskManager(taskSubhub, executor)
	taskService := taskservice.NewTaskService(taskRepo, taskSubhub, taskManager)

	tradingStrategy := trading.NewTradingStrategy(taskService, posSubhub, positionService, strategySubHub, taskSubhub, rpcGroupService)
	tradingTaskRepo := repository.NewTradingRepository(db)
	tradingService := trading.NewTradingService(tradingStrategy, strategySubHub, tradingTaskRepo, taskService, rpcGroupService)
	deps.TradingService = tradingService

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

	tradingController := controller.NewStrategyController(tradingService, walletService, rpcGroupService)
	tradingHandler := http.StripPrefix("/api/trading", handler.NewTradingHandler(tradingController))

	taskController := controller.NewTaskController(taskService, walletService)
	taskHandler := http.StripPrefix("/api/tasks", handler.NewTaskHandler(taskController))

	positionController := controller.PositionController{PositionService: positionService}
	positionHandler := http.StripPrefix("/api/position", handler.NewPositionHandler(&positionController))
	dashboardHandler := http.StripPrefix("/api/dashboard", handler.NewDashboardHandler(tradingController, taskController))

	mux := http.NewServeMux()
	mux.Handle("/api/tasks/", taskHandler)
	mux.Handle("/api/trading/", tradingHandler)
	mux.Handle("/api/position/", positionHandler)
	mux.Handle("/api/wallet/", walletHandler)
	mux.Handle("/api/dashboard/", dashboardHandler)
	mux.Handle("/api/settings/", settingsHandler)
	mux.Handle("/api/rpc_groups/", rpcGroupHandler)

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
