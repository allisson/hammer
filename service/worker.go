package service

import (
	"hash/fnv"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/allisson/go-pglock"
	"github.com/allisson/hammer"
	"go.uber.org/zap"
)

// Worker is a implementation of hammer.WorkerService
type Worker struct {
	lock            pglock.Locker
	deliveryService hammer.DeliveryService
	wg              sync.WaitGroup
	run             bool
}

func (w *Worker) stringToInt(s string) int64 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0
	}
	return int64(h.Sum32())
}

func (w *Worker) unlock(lockID int64) {
	err := w.lock.Unlock(lockID)
	if err != nil {
		logger.Error("unlock-delivery", zap.Error(err))
	}
}

func (w *Worker) dispatch(deliveryID string) {
	defer w.wg.Done()

	// Get lock
	lockID := w.stringToInt(deliveryID)
	ok, err := w.lock.Lock(lockID)
	if err != nil {
		logger.Error("lock-delivery", zap.Error(err))
		return
	}
	if !ok {
		return
	}
	defer w.unlock(lockID)

	// Get delivery
	delivery, err := w.deliveryService.Find(deliveryID)
	if err != nil {
		logger.Error("delivery-service-find", zap.Error(err))
		return
	}

	// Check delivery
	if delivery.Status != hammer.DeliveryStatusPending || time.Now().UTC().Before(delivery.ScheduledAt) {
		return
	}

	// Create http client with timeout
	httpClient := &http.Client{Timeout: time.Duration(delivery.DeliveryAttemptTimeout) * time.Second}

	// Dispatch
	deliveryAttempt, err := w.deliveryService.Dispatch(&delivery, httpClient)
	if err != nil {
		logger.Error("delivery-service-dispatch", zap.Error(err))
		return
	}

	if delivery.Status == hammer.DeliveryStatusCompleted {
		logger.Info(
			"delivery-made",
			zap.String("delivery_id", delivery.ID),
			zap.String("delivery_attempt_id", deliveryAttempt.ID),
			zap.Int("response_status_code", deliveryAttempt.ResponseStatusCode),
			zap.Int("execution_duration", deliveryAttempt.ExecutionDuration),
		)
	} else {
		logger.Info(
			"delivery-fail",
			zap.String("delivery_id", delivery.ID),
			zap.String("delivery_attempt_id", deliveryAttempt.ID),
			zap.Int("response_status_code", deliveryAttempt.ResponseStatusCode),
			zap.Int("execution_duration", deliveryAttempt.ExecutionDuration),
			zap.String("error", deliveryAttempt.Error),
			zap.Int("attempts", delivery.DeliveryAttempts),
			zap.Int("max_delivery_attempts", delivery.MaxDeliveryAttempts),
		)
	}
}

// Run worker flow
func (w *Worker) Run() error {
	for w.run {
		deliveries, err := w.deliveryService.FindToDispatch(hammer.WorkerDefaultFetchLimit, 0)
		if err != nil {
			return err
		}

		if len(deliveries) == 0 {
			time.Sleep(time.Duration(hammer.WorkerDatabaseDelay) * time.Second)
			continue
		}

		// Increment wait group
		w.wg.Add(len(deliveries))

		// Shuffle deliveries
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(deliveries), func(i, j int) { deliveries[i], deliveries[j] = deliveries[j], deliveries[i] })

		for _, deliveryID := range deliveries {
			id := deliveryID
			go w.dispatch(id)
		}

		// Wait for goroutines to finish
		w.wg.Wait()

		// Sleep with delay of hammer.WorkerDatabaseDelay
		time.Sleep(time.Duration(hammer.WorkerDatabaseDelay) * time.Second)
	}

	return nil
}

// Stop worker flow
func (w *Worker) Stop() error {
	w.run = false
	return nil
}

// NewWorker returns a new Worker
func NewWorker(lock pglock.Locker, deliveryService hammer.DeliveryService) Worker {
	return Worker{
		lock:            lock,
		deliveryService: deliveryService,
		run:             true,
	}
}
