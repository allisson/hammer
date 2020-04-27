package service

import (
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessage(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedMessage := hammer.MakeTestMessage()
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(messageRepo)
		messageRepo.On("Find", mock.Anything).Return(expectedMessage, nil)

		message, err := messageService.Find(expectedMessage.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessage, message)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedMessages := []hammer.Message{hammer.MakeTestMessage()}
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(messageRepo)
		messageRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedMessages, nil)

		messages, err := messageService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessages, messages)
	})

	t.Run("Test Create", func(t *testing.T) {
		message := hammer.MakeTestMessage()
		message.ID = ""
		messageRepo := &mocks.MessageRepository{}
		messageService := NewMessage(messageRepo)
		messageRepo.On("Store", mock.Anything).Return(nil)

		err := messageService.Create(&message)
		assert.Nil(t, err)
		assert.NotEqual(t, "", message.ID)
	})
}
