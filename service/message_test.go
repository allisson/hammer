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
		messageService := NewMessage(topicRepo, messageRepo)
		messageRepo.On("Find", mock.Anything).Return(expectedMessage, nil)

		message, err := messageService.Find(expectedMessage.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessage, message)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedMessages := []hammer.Message{hammer.MakeTestMessage()}
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(topicRepo, messageRepo)
		messageRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedMessages, nil)

		messages, err := messageService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessages, messages)
	})

	t.Run("Test FindByTopic", func(t *testing.T) {
		expectedMessages := []hammer.Message{hammer.MakeTestMessage()}
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(topicRepo, messageRepo)
		messageRepo.On("FindByTopic", mock.Anything, mock.Anything, mock.Anything).Return(expectedMessages, nil)

		messages, err := messageService.FindByTopic("topic_id", 50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessages, messages)
	})

	t.Run("Test Create", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		message.ID = ""
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(topicRepo, messageRepo)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, nil)
		messageRepo.On("Store", mock.Anything).Return(nil)

		err := messageService.Create(&message)
		assert.Nil(t, err)
		assert.NotEqual(t, "", message.ID)
	})

	t.Run("Test Create with topic does not exists on repository", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		message.ID = ""
		topicRepo := &mocks.TopicRepository{}
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(topicRepo, messageRepo)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, sql.ErrNoRows)

		err := messageService.Create(&message)
		assert.Equal(t, hammer.ErrTopicDoesNotExists, err)
	})
}
