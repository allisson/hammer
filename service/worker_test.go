package service

import (
	"testing"
	"time"

	lockmock "github.com/allisson/go-pglock/mocks"
	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorker(t *testing.T) {
	delivery := hammer.MakeTestDelivery()
	delivery.TopicID = "topic"
	delivery.SubscriptionID = "subscription"
	delivery.MessageID = "message"
	delivery.Status = hammer.DeliveryStatusCompleted
	deliveryAttempt := hammer.MakeTestDeliveryAttempt()
	deliveryAttempt.DeliveryID = delivery.ID
	deliveryService := &mocks.DeliveryService{}
	lock := &lockmock.Locker{}
	workerService := NewWorker(lock, deliveryService)
	deliveryService.On("FindToDispatch", hammer.WorkerDefaultFetchLimit, 0).Return([]string{delivery.ID}, nil)
	lock.On("Lock", mock.Anything).Return(true, nil)
	deliveryService.On("Find", delivery.ID).Return(delivery, nil)
	deliveryService.On("Dispatch", &delivery, mock.Anything).Return(deliveryAttempt, nil)
	lock.On("Unlock", mock.Anything).Return(nil)

	// Execute Run method in goroutine
	go func() {
		err := workerService.Run()
		assert.Nil(t, err)
	}()

	// Wait to Run method execute
	time.Sleep(10 * time.Millisecond)

	// Stop worker
	err := workerService.Stop()
	assert.Nil(t, err)
}
