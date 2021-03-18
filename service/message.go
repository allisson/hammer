package service

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"time"

	"github.com/allisson/hammer"
)

// Message is a implementation of hammer.MessageService
type Message struct {
	topicRepo        hammer.TopicRepository
	messageRepo      hammer.MessageRepository
	subscriptionRepo hammer.SubscriptionRepository
	deliveryRepo     hammer.DeliveryRepository
}

// Find returns hammer.Message by id
func (m Message) Find(ctx context.Context, id string) (*hammer.Message, error) {
	return m.messageRepo.Find(ctx, id)
}

// FindAll returns []hammer.Message by findOptions
func (m Message) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Message, error) {
	return m.messageRepo.FindAll(ctx, findOptions)
}

// Create a hammer.Message on repository
func (m Message) Create(ctx context.Context, message *hammer.Message) error {
	// Verify if topic already exists
	_, err := m.topicRepo.Find(ctx, message.TopicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrTopicDoesNotExists
		}
		return err
	}

	// Create message
	id, err := hammer.GenerateULID()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	message.ID = id
	message.CreatedAt = now
	message.Data = b64.StdEncoding.EncodeToString([]byte(message.Data))
	return m.messageRepo.Store(ctx, message)
}

// Delete a hammer.Message on repository
func (m Message) Delete(ctx context.Context, id string) error {
	return m.messageRepo.Delete(ctx, id)
}

// NewMessage returns a new Message with MessageRepo
func NewMessage(topicRepo hammer.TopicRepository, messageRepo hammer.MessageRepository, subscriptionRepo hammer.SubscriptionRepository, deliveryRepo hammer.DeliveryRepository) *Message {
	return &Message{
		topicRepo:        topicRepo,
		messageRepo:      messageRepo,
		subscriptionRepo: subscriptionRepo,
		deliveryRepo:     deliveryRepo,
	}
}
