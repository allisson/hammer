package service

import (
	"database/sql"
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessage(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedMessage := hammer.MakeTestMessage()
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo, txFactoryRepo)
		messageRepo.On("Find", mock.Anything).Return(expectedMessage, nil)

		message, err := messageService.Find(expectedMessage.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessage, message)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedMessages := []hammer.Message{hammer.MakeTestMessage()}
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo, txFactoryRepo)
		messageRepo.On("FindAll", mock.Anything).Return(expectedMessages, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		messages, err := messageService.FindAll(findOptions)
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
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo, txFactoryRepo)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, nil)
		subscriptionRepo.On("FindAll", mock.Anything).Return([]hammer.Subscription{hammer.MakeTestSubscription()}, nil)
		txFactoryRepo.On("New").Return(txRepo, nil)
		messageRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		deliveryRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		txRepo.On("Commit").Return(nil)

		err := messageService.Create(&message)
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
		txFactoryRepo := &mocks.TxFactoryRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo, txFactoryRepo)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, sql.ErrNoRows)

		err := messageService.Create(&message)
		assert.Equal(t, hammer.ErrTopicDoesNotExists, err)
	})

	t.Run("Test Delete", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		subscriptionRepo := &mocks.SubscriptionRepository{}
		deliveryRepo := &mocks.DeliveryRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		messageService := NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo, txFactoryRepo)
		messageRepo.On("Find", mock.Anything).Return(message, nil)
		txFactoryRepo.On("New").Return(txRepo, nil)
		messageRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
		txRepo.On("Commit").Return(nil)

		err := messageService.Delete(message.ID)
		assert.Nil(t, err)
	})
}
