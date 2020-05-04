package main

import (
	"hash/fnv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/allisson/go-env"
	"github.com/allisson/go-pglock"
	"github.com/allisson/hammer"
	repository "github.com/allisson/hammer/repository/postgres"
	"github.com/allisson/hammer/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sqlDB  *sqlx.DB
)

func stringToInt(s string) int64 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0
	}
	return int64(h.Sum32())
}

func init() {
	// Set logger
	logger, _ = zap.NewProduction()

	// Set database connection
	db, err := sqlx.Open("postgres", env.GetString("DATABASE_URL", ""))
	if err != nil {
		logger.Fatal("sqlx-open", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal("sqlx-ping", zap.Error(err))
	}
	sqlDB = db
}

func dispatch(lock *pglock.Lock, deliveryService hammer.DeliveryService, delivery *hammer.Delivery) {
	// Get lock
	lockID := stringToInt(delivery.ID)
	ok, err := lock.Lock(lockID)
	if err != nil {
		logger.Error("lock-delivery", zap.Error(err))
		return
	}
	if !ok {
		return
	}

	// Create http client with timeout
	httpClient := &http.Client{Timeout: time.Duration(delivery.DeliveryAttemptTimeout) * time.Second}

	// Dispatch
	err = deliveryService.Dispatch(delivery, httpClient)
	if err != nil {
		logger.Error("delivery-service-dispatch", zap.Error(err))
		return
	}
	logger.Info("delivery-made", zap.String("delivery_id", delivery.ID))

	// Unlock
	err = lock.Unlock(lockID)
	if err != nil {
		logger.Error("unlock-delivery", zap.Error(err))
	}
}

func getDeliveries(lock *pglock.Lock, deliveryService hammer.DeliveryService) {
	for {
		deliveries, err := deliveryService.FindToDispatch(hammer.WorkerDefaultFetchLimit, 0)
		if err != nil {
			logger.Error("get-deliveries-find-to-dispatch", zap.Error(err))
			continue
		}

		if len(deliveries) == 0 {
			continue
		}

		// Create wait group
		for _, delivery := range deliveries {
			dispatch(lock, deliveryService, &delivery)
		}

		// Sleep with delay of hammer.WorkerDatabaseDelay
		time.Sleep(time.Duration(hammer.WorkerDatabaseDelay) * time.Second)
	}
}

func main() {
	// Create a new lock
	lock := pglock.NewLock(sqlDB.DB)

	// Create repositories
	deliveryRepo := repository.NewDelivery(sqlDB)
	deliveryAttemptRepo := repository.NewDeliveryAttempt(sqlDB)
	txFactoryRepo := repository.NewTxFactory(sqlDB)

	// Create services
	deliveryService := service.NewDelivery(&deliveryRepo, &deliveryAttemptRepo, &txFactoryRepo)

	// Start dispatch deliveries
	go getDeliveries(&lock, &deliveryService)

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
		close(idleConnsClosed)
		logger.Info("worker-shutdown-finished")
	}()

	<-idleConnsClosed
}
