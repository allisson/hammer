package service

import (
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
		deliveryService := NewDelivery(deliveryRepo)
		deliveryRepo.On("Find", mock.Anything).Return(expectedDelivery, nil)

		delivery, err := deliveryService.Find(expectedDelivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDelivery, delivery)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedDeliveries := []hammer.Delivery{hammer.MakeTestDelivery()}
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryService := NewDelivery(deliveryRepo)
		deliveryRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedDeliveries, nil)

		deliveries, err := deliveryService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})

	t.Run("Test FindToDispatch", func(t *testing.T) {
		expectedDeliveries := []hammer.Delivery{hammer.MakeTestDelivery()}
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryService := NewDelivery(deliveryRepo)
		deliveryRepo.On("FindToDispatch", mock.Anything, mock.Anything).Return(expectedDeliveries, nil)

		deliveries, err := deliveryService.FindToDispatch(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})
}
