package service

import (
	"database/sql"
	"time"

	"github.com/allisson/hammer"
)

// Message is a implementation of hammer.MessageService
type Message struct {
	topicRepo        hammer.TopicRepository
	messageRepo      hammer.MessageRepository
	subscriptionRepo hammer.SubscriptionRepository
	deliveryRepo     hammer.DeliveryRepository
	txFactoryRepo    hammer.TxFactoryRepository
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

	// Start tx
	tx, err := m.txFactoryRepo.New()
	if err != nil {
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
	err = m.messageRepo.Store(tx, message)
	if err != nil {
		return err
	}

	// Get subscriptions
	subscriptions, err := m.subscriptionRepo.FindByTopic(message.TopicID)
	if err != nil {
		rollback(tx, "message-get-subscriptions")
		return err
	}

	// Create deliveries
	for _, subscription := range subscriptions {
		id, err := generateID()
		if err != nil {
			rollback(tx, "message-subscription-generate-id")
			return err
		}
		now := time.Now().UTC()
		delivery := hammer.Delivery{
			ID:                     id,
			TopicID:                message.TopicID,
			SubscriptionID:         subscription.ID,
			MessageID:              message.ID,
			Data:                   message.Data,
			URL:                    subscription.URL,
			SecretToken:            subscription.SecretToken,
			MaxDeliveryAttempts:    subscription.MaxDeliveryAttempts,
			DeliveryAttemptDelay:   subscription.DeliveryAttemptDelay,
			DeliveryAttemptTimeout: subscription.DeliveryAttemptTimeout,
			ScheduledAt:            now,
			Status:                 hammer.DeliveryStatusPending,
			CreatedAt:              now,
			UpdatedAt:              now,
		}
		err = m.deliveryRepo.Store(tx, &delivery)
		if err != nil {
			rollback(tx, "message-delivery-create-rollback")
			return err
		}
	}

	// tx Commit
	err = tx.Commit()
	if err != nil {
		rollback(tx, "message-create-rollback")
		return err
	}

	return nil
}

// NewMessage returns a new Message with MessageRepo
func NewMessage(topicRepo hammer.TopicRepository, messageRepo hammer.MessageRepository, subscriptionRepo hammer.SubscriptionRepository, deliveryRepo hammer.DeliveryRepository, txFactoryRepo hammer.TxFactoryRepository) Message {
	return Message{
		topicRepo:        topicRepo,
		messageRepo:      messageRepo,
		subscriptionRepo: subscriptionRepo,
		deliveryRepo:     deliveryRepo,
		txFactoryRepo:    txFactoryRepo,
	}
}
