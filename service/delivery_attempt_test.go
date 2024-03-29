package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
)

func TestDeliveryAttempt(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Find", func(t *testing.T) {
		expectedDeliveryAttempt := hammer.MakeTestDeliveryAttempt()
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		deliveryAttemptService := NewDeliveryAttempt(deliveryAttemptRepo)
		deliveryAttemptRepo.On("Find", mock.Anything, mock.Anything).Return(expectedDeliveryAttempt, nil)

		deliveryAttempt, err := deliveryAttemptService.Find(ctx, expectedDeliveryAttempt.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveryAttempt, deliveryAttempt)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedDeliveries := []*hammer.DeliveryAttempt{hammer.MakeTestDeliveryAttempt()}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		deliveryAttemptService := NewDeliveryAttempt(deliveryAttemptRepo)
		deliveryAttemptRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedDeliveries, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		deliveries, err := deliveryAttemptService.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})
}
