package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
)

func TestWorker(t *testing.T) {
	ctx := context.Background()
	delivery := hammer.MakeTestDelivery()
	delivery.TopicID = "topic"
	delivery.SubscriptionID = "subscription"
	delivery.MessageID = "message"
	delivery.Status = hammer.DeliveryStatusCompleted
	deliveryAttempt := hammer.MakeTestDeliveryAttempt()
	deliveryAttempt.DeliveryID = delivery.ID
	deliveryRepo := &mocks.DeliveryRepository{}
	workerService := NewWorker(deliveryRepo, 1*time.Second)

	deliveryRepo.On("Dispatch", mock.Anything).Return(&hammer.DeliveryAttempt{}, nil)

	// Execute Run method in goroutine
	go func() {
		workerService.Run(ctx)
	}()

	// Wait to Run method execute
	time.Sleep(10 * time.Millisecond)

	// Stop worker
	workerService.Stop(ctx)
}
