package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
)

func TestDelivery(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Find", func(t *testing.T) {
		expectedDelivery := hammer.MakeTestDelivery()
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo)
		deliveryRepo.On("Find", mock.Anything, mock.Anything).Return(expectedDelivery, nil)

		delivery, err := deliveryService.Find(ctx, expectedDelivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDelivery, delivery)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedDeliveries := []*hammer.Delivery{hammer.MakeTestDelivery()}
		deliveryRepo := &mocks.DeliveryRepository{}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		deliveryService := NewDelivery(deliveryRepo, deliveryAttemptRepo)
		deliveryRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedDeliveries, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		deliveries, err := deliveryService.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})
}
