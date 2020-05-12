package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDelivery(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedDelivery := hammer.MakeTestDelivery()
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo, txFactoryRepo)
		deliveryRepo.On("Find", mock.Anything).Return(expectedDelivery, nil)

		delivery, err := deliveryService.Find(expectedDelivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDelivery, delivery)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedDeliveries := []hammer.Delivery{hammer.MakeTestDelivery()}
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo, txFactoryRepo)
		deliveryRepo.On("FindAll", mock.Anything).Return(expectedDeliveries, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		deliveries, err := deliveryService.FindAll(findOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})

	t.Run("Test FindToDispatch", func(t *testing.T) {
		expectedDeliveries := []string{hammer.MakeTestDelivery().ID}
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo, txFactoryRepo)
		deliveryRepo.On("FindToDispatch", mock.Anything, mock.Anything).Return(expectedDeliveries, nil)

		deliveries, err := deliveryService.FindToDispatch(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})

	t.Run("Test Dispatch Success", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// nolint
			w.Write([]byte(`OK`))
		}))
		defer httpServer.Close()
		delivery := hammer.MakeTestDelivery()
		delivery.URL = httpServer.URL
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo, txFactoryRepo)
		txFactoryRepo.On("New").Return(txRepo, nil)
		deliveryAttemptRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		deliveryRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		txRepo.On("Commit").Return(nil)

		deliveryAttemps := delivery.DeliveryAttempts
		err := deliveryService.Dispatch(&delivery, httpServer.Client())
		assert.Nil(t, err)
		assert.Equal(t, delivery.DeliveryAttempts, deliveryAttemps+1)
		assert.Equal(t, hammer.DeliveryStatusCompleted, delivery.Status)
	})

	t.Run("Test Dispatch Error", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not_found", http.StatusNotFound)
		}))
		defer httpServer.Close()
		delivery := hammer.MakeTestDelivery()
		delivery.URL = httpServer.URL
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo, txFactoryRepo)
		txFactoryRepo.On("New").Return(txRepo, nil)
		deliveryAttemptRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		deliveryRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		txRepo.On("Commit").Return(nil)

		deliveryScheduledAt := delivery.ScheduledAt
		deliveryAttemps := delivery.DeliveryAttempts
		err := deliveryService.Dispatch(&delivery, httpServer.Client())
		assert.Nil(t, err)
		assert.Equal(t, delivery.DeliveryAttempts, deliveryAttemps+1)
		assert.Equal(t, hammer.DeliveryStatusPending, delivery.Status)
		assert.True(t, delivery.ScheduledAt.After(deliveryScheduledAt))
	})

	t.Run("Test Dispatch MaxDeliveryAttempts Error", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not_found", http.StatusNotFound)
		}))
		defer httpServer.Close()
		delivery := hammer.MakeTestDelivery()
		delivery.URL = httpServer.URL
		delivery.DeliveryAttempts = delivery.MaxDeliveryAttempts - 1
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo, txFactoryRepo)
		txFactoryRepo.On("New").Return(txRepo, nil)
		deliveryAttemptRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		deliveryRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		txRepo.On("Commit").Return(nil)

		deliveryScheduledAt := delivery.ScheduledAt
		deliveryAttemps := delivery.DeliveryAttempts
		err := deliveryService.Dispatch(&delivery, httpServer.Client())
		assert.Nil(t, err)
		assert.Equal(t, delivery.DeliveryAttempts, deliveryAttemps+1)
		assert.Equal(t, hammer.DeliveryStatusFailed, delivery.Status)
		assert.Equal(t, deliveryScheduledAt, delivery.ScheduledAt)
	})
}
