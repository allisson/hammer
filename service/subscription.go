package service

import (
	"database/sql"
	"time"

	"github.com/allisson/hammer"
)

// Subscription is a implementation of hammer.SubscriptionService
type Subscription struct {
	topicRepo        hammer.TopicRepository
	subscriptionRepo hammer.SubscriptionRepository
}

func (s *Subscription) topicExists(topicID string) error {
	_, err := s.topicRepo.Find(topicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrTopicDoesNotExists
		}
		return err
	}
	return nil
}

// Find returns hammer.Subscription by id
func (s *Subscription) Find(id string) (hammer.Subscription, error) {
	return s.subscriptionRepo.Find(id)
}

// FindAll returns []hammer.Subscription by limit and offset
func (s *Subscription) FindAll(limit, offset int) ([]hammer.Subscription, error) {
	return s.subscriptionRepo.FindAll(limit, offset)
}

// Create a hammer.Subscription on repository
func (s *Subscription) Create(subscription *hammer.Subscription) error {
	// Verify if subscription already exists
	_, err := s.subscriptionRepo.Find(subscription.ID)
	if err == nil {
		return hammer.ErrSubscriptionAlreadyExists
	}

	// Verify if topic already exists
	err = s.topicExists(subscription.TopicID)
	if err != nil {
		return err
	}

	// Create new subscription with default values
	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	if subscription.SecretToken == "" {
		id, err := generateID()
		if err != nil {
			return err
		}
		subscription.SecretToken = id
	}
	if subscription.MaxDeliveryAttempts == 0 {
		subscription.MaxDeliveryAttempts = 5
	}
	if subscription.DeliveryAttemptDelay == 0 {
		subscription.DeliveryAttemptDelay = 60
	}
	if subscription.DeliveryAttemptTimeout == 0 {
		subscription.DeliveryAttemptTimeout = 5
	}
	return s.subscriptionRepo.Store(subscription)
}

// Update a hammer.Subscription on repository
func (s *Subscription) Update(subscription *hammer.Subscription) error {
	// Verify if subscription already exists
	subscriptionFromRepo, err := s.subscriptionRepo.Find(subscription.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrSubscriptionDoesNotExists
		}
		return err
	}

	// Verify if topic already exists
	err = s.topicExists(subscription.TopicID)
	if err != nil {
		return err
	}

	// Update subscription
	subscription.ID = subscriptionFromRepo.ID
	subscription.TopicID = subscriptionFromRepo.TopicID
	subscription.UpdatedAt = time.Now().UTC()
	return s.subscriptionRepo.Store(subscription)
}

// NewSubscription returns a new Subscription with SubscriptionRepo
func NewSubscription(topicRepo hammer.TopicRepository, subscriptionRepo hammer.SubscriptionRepository) Subscription {
	return Subscription{
		topicRepo:        topicRepo,
		subscriptionRepo: subscriptionRepo,
	}
}
