package main

import (
	"hash/fnv"
	"math/rand"
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

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

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

func (t *taskJob) unlock(lockID int64) {
	err := t.lock.Unlock(lockID)
	if err != nil {
		logger.Error("unlock-delivery", zap.Error(err))
	}
}

func (t *taskJob) DeliveriesToDispatch() ([]string, error) {
	return t.deliveryService.FindToDispatch(hammer.WorkerDefaultFetchLimit, 0)
}

func (t *taskJob) Dispatch(deliveryID string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get lock
	lockID := stringToInt(deliveryID)
	ok, err := t.lock.Lock(lockID)
	if err != nil {
		logger.Error("lock-delivery", zap.Error(err))
		return
	}
	if !ok {
		return
	}
	defer t.unlock(lockID)

	// Get delivery
	delivery, err := t.deliveryService.Find(deliveryID)
	if err != nil {
		return
	}

	// Check delivery
	if delivery.Status != hammer.DeliveryStatusPending || delivery.ScheduledAt.After(time.Now().UTC()) {
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

		// Shuffle deliveries
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(deliveries), func(i, j int) { deliveries[i], deliveries[j] = deliveries[j], deliveries[i] })

		// Create wait group
		for _, deliveryID := range deliveries {
			// Random sleep to avoid lock race condition
			time.Sleep(time.Duration(randomInt(100, 200)) * time.Millisecond)
			// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			go func(deliveryID string, wg *sync.WaitGroup) {
				job.Dispatch(deliveryID, wg)
			}(deliveryID, &wg)
		}

		// Wait for goroutines to finish
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
