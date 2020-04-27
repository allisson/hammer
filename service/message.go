package service

import (
	"time"

	"github.com/allisson/hammer"
)

// Message is a implementation of hammer.MessageService
type Message struct {
	messageRepo hammer.MessageRepository
}

// Find returns hammer.Message by id
func (m *Message) Find(id string) (hammer.Message, error) {
	return m.messageRepo.Find(id)
}

// FindAll returns []hammer.Message by limit and offset
func (m *Message) FindAll(limit, offset int) ([]hammer.Message, error) {
	return m.messageRepo.FindAll(limit, offset)
}

// Create a hammer.Message on repository
func (m *Message) Create(message *hammer.Message) error {
	id, err := generateID()
	if err != nil {
		return err
	}
	message.ID = id
	message.CreatedAt = time.Now().UTC()
	message.CreatedDeliveries = false
	return m.messageRepo.Store(message)
}

// NewMessage returns a new Message with MessageRepo
func NewMessage(messageRepo hammer.MessageRepository) Message {
	return Message{messageRepo: messageRepo}
}
