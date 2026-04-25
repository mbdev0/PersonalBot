package main

import (
	"context"
	"database/sql"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"personal_bot/app/wire"
	"personal_bot/infrastructure/persistence"
	"syscall"
	"time"

	"personal_bot/pkg/logger"
)

func main() {
	db, err := persistence.NewConnection()
	if err != nil {
		panic(err)
	}

	server, shutdownServices := wire.WireServicesAndLaunchServer(db)

	go func() {
		logger.Information("Starting server on port 9090:")
		logger.Information("http://localhost:9090")
		if err := server.ListenAndServe(); err != nil {
			logger.Information("stopped listening to localhost:9090")
		}
	}()

	// checkGoRoutines()
	handleShutdown(db, server, shutdownServices)
}

func handleShutdown(db *sql.DB, server *http.Server, shutdownServices wire.ShutdownServices) {
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdown.Done()

	logger.Information("shutting down application")

	logger.Information("pushing tasks to db")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := shutdownServices.TradingService.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}

	// do the tasks here
	err = shutdownServices.TaskService.Shutdown(ctx)
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
