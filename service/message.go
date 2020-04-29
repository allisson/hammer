package service

import (
	"database/sql"
	"time"

	"github.com/allisson/hammer"
)

// Message is a implementation of hammer.MessageService
type Message struct {
	topicRepo   hammer.TopicRepository
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
	// Verify if topic already exists
	_, err := m.topicRepo.Find(message.TopicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrTopicDoesNotExists
		}
		return err
	}

	// Create message
	id, err := generateID()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	message.ID = id
	message.CreatedAt = now
	message.UpdatedAt = now
	message.CreatedDeliveries = false
	return m.messageRepo.Store(message)
}

// NewMessage returns a new Message with MessageRepo
func NewMessage(topicRepo hammer.TopicRepository, messageRepo hammer.MessageRepository) Message {
	return Message{
		topicRepo:   topicRepo,
		messageRepo: messageRepo,
	}
}
