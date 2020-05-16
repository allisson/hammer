package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/allisson/go-env"
	"github.com/allisson/go-pglock"
	repository "github.com/allisson/hammer/repository/postgres"
	"github.com/allisson/hammer/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	logger  *zap.Logger
	sqlDB   *sqlx.DB
	sqlConn *sql.Conn
)

func init() {
	// Set logger
	logger, _ = zap.NewProduction()

	// Set database connection
	db, err := sqlx.Open("postgres", env.GetString("HAMMER_DATABASE_URL", ""))
	if err != nil {
		logger.Fatal("sqlx-open", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal("sqlx-ping", zap.Error(err))
	}
	sqlDB = db
	conn, err := sqlDB.DB.Conn(context.Background())
	if err != nil {
		logger.Fatal("sql-conn", zap.Error(err))
	}
	sqlConn = conn
}

func main() {
	// Create a new lock
	lock := pglock.NewLock(sqlConn)

	// Create repositories
	deliveryRepo := repository.NewDelivery(sqlDB)
	deliveryAttemptRepo := repository.NewDeliveryAttempt(sqlDB)
	txFactoryRepo := repository.NewTxFactory(sqlDB)

	// Create services
	deliveryService := service.NewDelivery(&deliveryRepo, &deliveryAttemptRepo, &txFactoryRepo)
	workerService := service.NewWorker(&lock, &deliveryService)

	// Make graceful shutdown
	logger.Info("worker-started")
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		logger.Info("worker-shutdown-started")
		if err := workerService.Stop(); err != nil {
			logger.Error("worker-service-stop", zap.Error(err))
		}
		close(idleConnsClosed)
		logger.Info("worker-shutdown-finished")
	}()

	if err := workerService.Run(); err != nil {
		logger.Error("worker-service-run", zap.Error(err))
	}

	<-idleConnsClosed
}
