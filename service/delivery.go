package service

import (
	"time"

	"github.com/allisson/hammer"
)

// Delivery is a implementation of hammer.DeliveryService
type Delivery struct {
	deliveryRepo     hammer.DeliveryRepository
	subscriptionRepo hammer.SubscriptionRepository
	messageRepo      hammer.MessageRepository
}

// Find returns hammer.Delivery by id
func (d *Delivery) Find(id string) (hammer.Delivery, error) {
	return d.deliveryRepo.Find(id)
}

// FindAll returns []hammer.Delivery by limit and offset
func (d *Delivery) FindAll(limit, offset int) ([]hammer.Delivery, error) {
	return d.deliveryRepo.FindAll(limit, offset)
}

// Create a list of hammer.Delivery from message
func (d *Delivery) Create(message *hammer.Message) ([]hammer.Delivery, error) {
	deliveries := []hammer.Delivery{}

	// Get subscriptions
	subscriptions, err := d.subscriptionRepo.FindByTopic(message.TopicID)
	if err != nil {
		return deliveries, err
	}

	// Create deliveries
	for _, subscription := range subscriptions {
		id, err := generateID()
		if err != nil {
			return deliveries, err
		}
		now := time.Now().UTC()
		delivery := hammer.Delivery{
			ID:             id,
			TopicID:        message.TopicID,
			SubscriptionID: subscription.ID,
			MessageID:      message.ID,
			ScheduledAt:    now,
			Status:         "pending",
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		err = d.deliveryRepo.Store(&delivery)
		if err != nil {
			return deliveries, err
		}
		deliveries = append(deliveries, delivery)
	}

	// Update message CreatedDeliveries
	message.CreatedDeliveries = true
	err = d.messageRepo.Store(message)

	return deliveries, err
}

// NewDelivery returns a new Delivery with DeliveryRepo
func NewDelivery(deliveryRepo hammer.DeliveryRepository, subscriptionRepo hammer.SubscriptionRepository, messageRepo hammer.MessageRepository) Delivery {
	return Delivery{
		deliveryRepo:     deliveryRepo,
		subscriptionRepo: subscriptionRepo,
		messageRepo:      messageRepo,
	}
}
