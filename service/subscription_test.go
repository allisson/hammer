package service

import (
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscription(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedSubscription := hammer.MakeTestSubscription()
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(subscriptionRepo)
		subscriptionRepo.On("Find", mock.Anything).Return(expectedSubscription, nil)

		subscription, err := subscriptionService.Find(expectedSubscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedSubscription, subscription)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedSubscriptions := []hammer.Subscription{hammer.MakeTestSubscription()}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(subscriptionRepo)
		subscriptionRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedSubscriptions, nil)

		subscriptions, err := subscriptionService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedSubscriptions, subscriptions)
	})

	t.Run("Test Create", func(t *testing.T) {
		subscription := hammer.MakeTestSubscription()
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(subscriptionRepo)
		subscriptionRepo.On("Store", mock.Anything).Return(nil)

		err := subscriptionService.Create(&subscription)
		assert.Nil(t, err)
	})

	t.Run("Test Update", func(t *testing.T) {
		subscription := hammer.MakeTestSubscription()
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(subscriptionRepo)
		subscriptionRepo.On("Store", mock.Anything).Return(nil)

		err := subscriptionService.Create(&subscription)
		assert.Nil(t, err)

		subscription.Name = "My Subscription"
		err = subscriptionService.Update(&subscription)
		assert.Nil(t, err)
	})
}
