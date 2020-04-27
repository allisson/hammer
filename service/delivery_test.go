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
		subscriptionRepo := &mocks.SubscriptionRepository{}
		messageRepo := &mocks.MessageRepository{}
		deliveryService := NewDelivery(deliveryRepo, subscriptionRepo, messageRepo)
		deliveryRepo.On("Find", mock.Anything).Return(expectedDelivery, nil)

		delivery, err := deliveryService.Find(expectedDelivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDelivery, delivery)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedDeliverys := []hammer.Delivery{hammer.MakeTestDelivery()}
		deliveryRepo := &mocks.DeliveryRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		messageRepo := &mocks.MessageRepository{}
		deliveryService := NewDelivery(deliveryRepo, subscriptionRepo, messageRepo)
		deliveryRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedDeliverys, nil)

		deliveries, err := deliveryService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliverys, deliveries)
	})

	t.Run("Test Create", func(t *testing.T) {
		subscription1 := hammer.MakeTestSubscription()
		subscription2 := hammer.MakeTestSubscription()
		message := hammer.MakeTestMessage()
		deliveryRepo := &mocks.DeliveryRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		messageRepo := &mocks.MessageRepository{}
		deliveryService := NewDelivery(deliveryRepo, subscriptionRepo, messageRepo)
		subscriptionRepo.On("FindByTopic", mock.Anything).Return([]hammer.Subscription{subscription1, subscription2}, nil)
		deliveryRepo.On("Store", mock.Anything).Return(nil)
		messageRepo.On("Store", mock.Anything).Return(nil)

		deliveries, err := deliveryService.Create(&message)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(deliveries))
		assert.Equal(t, true, message.CreatedDeliveries)
	})
}
