package server

import (
	"context"
	"database/sql"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"personal_bot/backend/app/wire"
	"personal_bot/backend/infrastructure/persistence"
	"syscall"
	"time"

	"personal_bot/backend/pkg/logger"
)

type Server struct {
	http     *http.Server
	db       *sql.DB
	shutdown wire.ShutdownServices
}

func New(ctx context.Context) (*Server, error) {
	db, err := persistence.NewConnection()
	if err != nil {
		return nil, err
	}

	httpServer, shutdownServices := wire.WireServicesAndLaunchServer(ctx, db)

	return &Server{
		http:     httpServer,
		db:       db,
		shutdown: shutdownServices,
	}, nil
}

func (s *Server) Start() {
	logger.Information("Starting server on port 9090")
	if err := s.http.ListenAndServe(); err != nil {
		logger.Error(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.shutdown.TradingService.Shutdown(ctx)
	s.shutdown.TaskService.Shutdown(ctx)
	s.http.Shutdown(ctx)
	s.db.Close()
}

func main() {
	mainCtx := context.Background()

	db, err := persistence.NewConnection()
	if err != nil {
		panic(err)
	}

	server, shutdownServices := wire.WireServicesAndLaunchServer(mainCtx, db)

	go func() {
		logger.Information("Starting server on port 9090:")
		logger.Information("http://localhost:9090")
		if err := server.ListenAndServe(); err != nil {
			logger.Error(err)
			logger.Information("stopped listening to localhost:9090")
		}
	}()

	// checkGoRoutines()
	handleShutdown(mainCtx, db, server, shutdownServices)
}

func handleShutdown(ctx context.Context, db *sql.DB, server *http.Server, shutdownServices wire.ShutdownServices) {
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdown.Done()

	logger.Information("shutting down application")

	logger.Information("pushing tasks to db")
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
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
