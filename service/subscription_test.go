package service

import (
	"database/sql"
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscription(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedSubscription := hammer.MakeTestSubscription()
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		subscriptionRepo.On("Find", mock.Anything).Return(expectedSubscription, nil)

		subscription, err := subscriptionService.Find(expectedSubscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedSubscription, subscription)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedSubscriptions := []hammer.Subscription{hammer.MakeTestSubscription()}
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		subscriptionRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedSubscriptions, nil)

		subscriptions, err := subscriptionService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedSubscriptions, subscriptions)
	})

	t.Run("Test Create", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, nil)
		subscriptionRepo.On("Store", mock.Anything).Return(nil)
		subscriptionRepo.On("Find", mock.Anything).Return(hammer.Subscription{}, sql.ErrNoRows)

		err := subscriptionService.Create(&subscription)
		assert.Nil(t, err)
		assert.NotEqual(t, "", subscription.SecretToken)
	})

	t.Run("Test Create with topic does not exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, sql.ErrNoRows)
		subscriptionRepo.On("Find", mock.Anything).Return(hammer.Subscription{}, sql.ErrNoRows)

		err := subscriptionService.Create(&subscription)
		assert.Equal(t, hammer.ErrTopicDoesNotExists, err)
	})

	t.Run("Test Create with subscription already exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, nil)
		subscriptionRepo.On("Find", mock.Anything).Return(hammer.Subscription{}, nil)

		err := subscriptionService.Create(&subscription)
		assert.Equal(t, hammer.ErrSubscriptionAlreadyExists, err)
	})

	t.Run("Test Update", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, nil)
		subscriptionRepo.On("Store", mock.Anything).Return(nil)
		subscriptionRepo.On("Find", mock.Anything).Return(hammer.Subscription{}, nil)

		subscription.Name = "My Subscription"
		err := subscriptionService.Update(&subscription)
		assert.Nil(t, err)
	})

	t.Run("Test Update with topic does not exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, sql.ErrNoRows)
		subscriptionRepo.On("Find", mock.Anything).Return(hammer.Subscription{}, sql.ErrNoRows)

		subscription.Name = "My Subscription"
		err := subscriptionService.Update(&subscription)
		assert.Equal(t, hammer.ErrSubscriptionDoesNotExists, err)
	})

	t.Run("Test Update with subscription does not exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		topicRepo := &mocks.TopicRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		subscriptionService := NewSubscription(topicRepo, subscriptionRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, nil)
		subscriptionRepo.On("Find", mock.Anything).Return(hammer.Subscription{}, sql.ErrNoRows)

		subscription.Name = "My Subscription"
		err := subscriptionService.Update(&subscription)
		assert.Equal(t, hammer.ErrSubscriptionDoesNotExists, err)
	})
}
