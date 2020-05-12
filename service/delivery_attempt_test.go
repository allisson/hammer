package service

import (
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeliveryAttempt(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedDeliveryAttempt := hammer.MakeTestDeliveryAttempt()
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		deliveryAttemptService := NewDeliveryAttempt(deliveryAttemptRepo)
		deliveryAttemptRepo.On("Find", mock.Anything).Return(expectedDeliveryAttempt, nil)

		deliveryAttempt, err := deliveryAttemptService.Find(expectedDeliveryAttempt.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveryAttempt, deliveryAttempt)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedDeliveries := []hammer.DeliveryAttempt{hammer.MakeTestDeliveryAttempt()}
		deliveryAttemptRepo := &mocks.DeliveryAttemptRepository{}
		deliveryAttemptService := NewDeliveryAttempt(deliveryAttemptRepo)
		deliveryAttemptRepo.On("FindAll", mock.Anything).Return(expectedDeliveries, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		deliveries, err := deliveryAttemptService.FindAll(findOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveries, deliveries)
	})
}
