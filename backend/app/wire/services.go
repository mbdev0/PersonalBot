package wire

import (
	"context"
	"database/sql"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/api/handler"
	"personal_bot/app"
	programnames "personal_bot/app/program_names"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/core/program"
	cryptostates "personal_bot/internal/services/crypto_states"
	"personal_bot/internal/services/notifier"
	"personal_bot/internal/services/position"
	rpcgroups "personal_bot/internal/services/rpc_groups"
	"personal_bot/internal/services/settings"
	positionhub "personal_bot/internal/services/subscription_hub/position"
	"personal_bot/internal/services/subscription_hub/strategy"
	"personal_bot/internal/services/subscription_hub/task"
	"personal_bot/internal/services/wallet"
	datastream "personal_bot/internal/solana/monitoring/data_stream"
	pumpfunamm "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm"
	pumpfunnative "personal_bot/internal/solana/programs/pumpfun/pumpfun_native"
	"personal_bot/pkg/logger"

	taskservice "personal_bot/internal/services/task_service"
	"personal_bot/internal/services/trading"
	"personal_bot/internal/services/transaction"

	"github.com/gagliardetto/solana-go/rpc"
)

type ShutdownServices struct {
	TaskService    *taskservice.TaskService
	TradingService *trading.Service
}

func WireServicesAndLaunchServer(mainCtx context.Context, db *sql.DB) (*http.Server, ShutdownServices) {

	settingsRepo := repository.NewSettingRepository(db)
	settingsService := settings.NewSettingsService(settingsRepo)
	settings, err := settingsService.GetSettings(mainCtx)
	if err != nil {
		logger.Error(err)
	}

	discordNotifier := notifier.NewDiscordNotifier(settings)
	settingsController := controller.NewSettingsController(settingsService, discordNotifier)

	taskSubhub := task.NewTaskSubscriptionHub()
	posSubhub := positionhub.NewSubscriptionHub(settingsService)
	posRepo := repository.NewPositionRepository(db)
	positionService := position.NewPositionService(posSubhub, posRepo)

	rpcGroupRepo := repository.NewRPCGroupRepository(db)
	rpcGroupService := rpcgroups.NewRPCGroupService(rpcGroupRepo)
	rpcGroupController := controller.NewRPCGroupController(rpcGroupService)

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
	deps.TaskService = taskService

	tradingStrategy := trading.NewTradingStrategy(taskService, posSubhub, positionService, strategySubHub, taskSubhub, rpcGroupService)
	tradingTaskRepo := repository.NewTradingRepository(db)
	tradingService := trading.NewTradingService(tradingStrategy, strategySubHub, tradingTaskRepo, taskService, rpcGroupService)
	deps.TradingService = tradingService

	walletRepo := repository.NewWalletRepository(db)
	walletService := wallet.NewWalletService(walletRepo)

	solanaStreamManager := datastream.NewSolanaStream(mainCtx)
	registerPrograms(solanaStreamManager)

	walletController := controller.NewWalletController(walletService)
	tradingController := controller.NewStrategyController(tradingService, walletService, rpcGroupService)
	taskController := controller.NewTaskController(taskService, walletService, rpcGroupService)
	positionController := controller.NewPositionController(positionService)

	tradingHandler := http.StripPrefix("/api/trading", handler.NewTradingHandler(tradingController))
	taskHandler := http.StripPrefix("/api/tasks", handler.NewTaskHandler(taskController))
	rpcGroupHandler := http.StripPrefix("/api/rpc_groups", handler.NewRPCGroupHandler(rpcGroupController))
	settingsHandler := http.StripPrefix("/api/settings", handler.NewSettingsHandler(settingsController))
	positionHandler := http.StripPrefix("/api/position", handler.NewPositionHandler(&positionController))
	dashboardHandler := http.StripPrefix("/api/dashboard", handler.NewDashboardHandler(tradingController, taskController))
	marketHandler := http.StripPrefix("/api/market", handler.NewMarketHandler())
	walletHandler := http.StripPrefix("/api/wallet", handler.NewWalletHandler(walletController))

	mux := http.NewServeMux()
	mux.Handle("/api/tasks/", taskHandler)
	mux.Handle("/api/trading/", tradingHandler)
	mux.Handle("/api/position/", positionHandler)
	mux.Handle("/api/wallet/", walletHandler)
	mux.Handle("/api/dashboard/", dashboardHandler)
	mux.Handle("/api/settings/", settingsHandler)
	mux.Handle("/api/rpc_groups/", rpcGroupHandler)
	mux.Handle("/api/market/", marketHandler)

	loadFromDB(mainCtx, tradingService, taskService, positionService)
	app.Launch(rpc.New(settings.PositionNodes.HTTPNode))

	server := &http.Server{
		Addr:    ":9090",
		Handler: CORS(mux.ServeHTTP),
	}

	shutdownServices := ShutdownServices{
		TaskService:    taskService,
		TradingService: tradingService,
	}

	return server, shutdownServices
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

func loadFromDB(ctx context.Context, tradingService *trading.Service, taskService *taskservice.TaskService, positionService *position.Service) {
	err := tradingService.LoadFromDB(ctx)
	if err != nil {
		logger.Error(err)
	}

	err = taskService.LoadFromDB(ctx)
	if err != nil {
		logger.Error("error whilst loading from db: ", err)
	}

	err = positionService.LoadFromDB(ctx)
	if err != nil {
		logger.Error("error whilst loading from db: ", err)
	}
}

func registerPrograms(solanaStream datastream.DataStream) {
	program.Register(string(programnames.PumpfunNative), pumpfunnative.New(solanaStream))
	program.Register(string(programnames.PumpfunAMM), pumpfunamm.New(solanaStream))
}
