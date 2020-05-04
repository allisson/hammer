package service

import (
	"database/sql"
	"time"

	"github.com/allisson/hammer"
	"go.uber.org/zap"
)

// Message is a implementation of hammer.MessageService
type Message struct {
	topicRepo     hammer.TopicRepository
	messageRepo   hammer.MessageRepository
	txFactoryRepo hammer.TxFactoryRepository
}

// Find returns hammer.Message by id
func (m *Message) Find(id string) (hammer.Message, error) {
	return m.messageRepo.Find(id)
}

// FindAll returns []hammer.Message by limit and offset
func (m *Message) FindAll(limit, offset int) ([]hammer.Message, error) {
	return m.messageRepo.FindAll(limit, offset)
}

// FindByTopic returns []hammer.Message by topicID, limit and offset
func (m *Message) FindByTopic(topicID string, limit, offset int) ([]hammer.Message, error) {
	return m.messageRepo.FindByTopic(topicID, limit, offset)
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
	tx, err := m.txFactoryRepo.New()
	if err != nil {
		return err
	}
	id, err := generateID()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	message.ID = id
	message.CreatedAt = now
	message.UpdatedAt = now
	message.CreatedDeliveries = false
	err = m.messageRepo.Store(tx, message)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			logger.Error("message-create-rollback", zap.Error(rErr))
		}
		return err
	}

	return nil
}

// NewMessage returns a new Message with MessageRepo
func NewMessage(topicRepo hammer.TopicRepository, messageRepo hammer.MessageRepository, txFactoryRepo hammer.TxFactoryRepository) Message {
	return Message{
		topicRepo:     topicRepo,
		messageRepo:   messageRepo,
		txFactoryRepo: txFactoryRepo,
	}
}
