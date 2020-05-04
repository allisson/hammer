package main

import (
	"hash/fnv"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

type taskJob struct {
	lock            *pglock.Lock
	deliveryService hammer.DeliveryService
}

func (t *taskJob) DeliveriesToDispatch() ([]hammer.Delivery, error) {
	return t.deliveryService.FindToDispatch(hammer.WorkerDefaultFetchLimit, 0)
}

func (t *taskJob) Dispatch(delivery hammer.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get lock
	lockID := stringToInt(delivery.ID)
	ok, err := t.lock.Lock(lockID)
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
	err = t.deliveryService.Dispatch(&delivery, httpClient)
	if err != nil {
		logger.Error("delivery-service-dispatch", zap.Error(err))
		return
	}
	logger.Info("delivery-made", zap.String("delivery_id", delivery.ID), zap.String("message_id", delivery.MessageID))

	// Unlock
	err = t.lock.Unlock(lockID)
	if err != nil {
		logger.Error("unlock-delivery", zap.Error(err))
	}
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

func getDeliveries(job *taskJob) {
	for {
		deliveries, err := job.DeliveriesToDispatch()
		if err != nil {
			logger.Error("get-deliveries-find-to-dispatch", zap.Error(err))
			continue
		}

		if len(deliveries) == 0 {
			continue
		}

		// Create wait group
		wg := sync.WaitGroup{}
		wg.Add(len(deliveries))

		// Create wait group
		for _, delivery := range deliveries {
			// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			go func(delivery hammer.Delivery, wg *sync.WaitGroup) {
				job.Dispatch(delivery, wg)
			}(delivery, &wg)
		}

		wg.Wait()

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

	// Create taskJob
	job := taskJob{lock: &lock, deliveryService: &deliveryService}

	// Start dispatch deliveries
	go getDeliveries(&job)

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
