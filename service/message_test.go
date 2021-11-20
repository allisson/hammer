package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
)

func TestMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Find", func(t *testing.T) {
		expectedMessage := hammer.MakeTestMessage()
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo)
		messageRepo.On("Find", mock.Anything, expectedMessage.ID).Return(expectedMessage, nil)

		message, err := messageService.Find(ctx, expectedMessage.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessage, message)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		expectedMessages := []*hammer.Message{message}
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo)
		messageRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedMessages, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		messages, err := messageService.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessages, messages)
	})

	t.Run("Test Create", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		message.ID = ""
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo)
		topicRepo.On("Find", mock.Anything, message.ID).Return(&hammer.Topic{}, nil)
		subscriptionRepo.On("FindAll", mock.Anything, mock.Anything).Return([]*hammer.Subscription{hammer.MakeTestSubscription()}, nil)
		messageRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		deliveryRepo.On("Store", mock.Anything, mock.Anything).Return(nil)

		err := messageService.Create(ctx, message)
		assert.Nil(t, err)
		assert.NotEqual(t, "", message.ID)
		assert.Equal(t, "eyJpZCI6ICJpZCIsICJuYW1lIjogIkFsbGlzc29uIn0=", message.Data)
	})

	t.Run("Test Create with topic does not exists on repository", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		message.ID = ""
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(&hammer.Topic{}, sql.ErrNoRows)

		err := messageService.Create(ctx, message)
		assert.Equal(t, hammer.ErrTopicDoesNotExists, err)
	})

	t.Run("Test Delete", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo)
		messageRepo.On("Find", mock.Anything).Return(message, nil)
		messageRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

		err := messageService.Delete(ctx, message.ID)
		assert.Nil(t, err)
	})
}
