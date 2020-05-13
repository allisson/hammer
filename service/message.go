package service

import (
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
	txFactoryRepo    hammer.TxFactoryRepository
}

// Find returns hammer.Message by id
func (m *Message) Find(id string) (hammer.Message, error) {
	return m.messageRepo.Find(id)
}

// FindAll returns []hammer.Message by findOptions
func (m *Message) FindAll(findOptions hammer.FindOptions) ([]hammer.Message, error) {
	return m.messageRepo.FindAll(findOptions)
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
	id, err := generateULID()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	message.ID = id
	message.CreatedAt = now
	message.Data = b64.StdEncoding.EncodeToString([]byte(message.Data))
	err = m.messageRepo.Store(tx, message)
	if err != nil {
		return err
	}

	// Get subscriptions
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "topic_id",
				Operator:  "=",
				Value:     message.TopicID,
			},
		},
	}
	subscriptions, err := m.subscriptionRepo.FindAll(findOptions)
	if err != nil {
		rollback(tx, "message-get-subscriptions")
		return err
	}

	// Create deliveries
	for _, subscription := range subscriptions {
		id, err := generateULID()
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
