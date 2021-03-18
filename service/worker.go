package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/allisson/hammer"
	"go.uber.org/zap"
)

// Worker is a implementation of hammer.WorkerService
type Worker struct {
	deliveryRepo    hammer.DeliveryRepository
	pollingInterval time.Duration
	isRun           bool
}

func (w *Worker) run(ctx context.Context) {
	for w.isRun {
		// Dispatch webhook.
		deliveryAttempt, err := w.deliveryRepo.Dispatch(ctx)
		if err != nil {
			zap.L().Error("worker-dispatch-error", zap.Error(err))
			time.Sleep(w.pollingInterval)
			continue
		}
		if deliveryAttempt == nil {
			time.Sleep(w.pollingInterval)
			continue
		}

		// Log delivery attempt.
		zap.L().Info(
			"worker-delivery-attempt-created",
			zap.String("id", deliveryAttempt.ID),
			zap.String("delivery_id", deliveryAttempt.DeliveryID),
			zap.Int("response_status_code", deliveryAttempt.ResponseStatusCode),
			zap.Int("execution_duration", deliveryAttempt.ExecutionDuration),
			zap.Bool("success", deliveryAttempt.Success),
		)
	}

	zap.L().Info("worker-shutdown-completed")
}

// Run worker flow
func (w *Worker) Run(ctx context.Context) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		w.Stop(ctx)
		close(idleConnsClosed)
	}()

	zap.L().Info("worker-started")
	w.run(ctx)

	<-idleConnsClosed
}

// Stop worker flow
func (w *Worker) Stop(ctx context.Context) {
	w.isRun = false
	zap.L().Info("worker-shutdown-started")
}

// NewWorker returns a new Worker
func NewWorker(deliveryRepo hammer.DeliveryRepository, pollingInterval time.Duration) Worker {
	return Worker{
		deliveryRepo:    deliveryRepo,
		pollingInterval: pollingInterval,
		isRun:           true,
	}
}
