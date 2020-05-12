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
	txFactoryRepo    hammer.TxFactoryRepository
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

// FindAll returns []hammer.Subscription by findOptions
func (s *Subscription) FindAll(findOptions hammer.FindOptions) ([]hammer.Subscription, error) {
	return s.subscriptionRepo.FindAll(findOptions)
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
	tx, err := s.txFactoryRepo.New()
	if err != nil {
		return err
	}
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
	if subscription.MaxDeliveryAttempts <= 0 {
		subscription.MaxDeliveryAttempts = hammer.DefaultMaxDeliveryAttempts
	}
	if subscription.DeliveryAttemptDelay <= 0 {
		subscription.DeliveryAttemptDelay = hammer.DefaultDeliveryAttemptDelay
	}
	if subscription.DeliveryAttemptTimeout <= 0 {
		subscription.DeliveryAttemptTimeout = hammer.DefaultDeliveryAttemptTimeout
	}
	err = s.subscriptionRepo.Store(tx, subscription)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		rollback(tx, "subscription-create-rollback")
		return err
	}

	return nil
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
	tx, err := s.txFactoryRepo.New()
	if err != nil {
		return err
	}
	subscription.ID = subscriptionFromRepo.ID
	subscription.TopicID = subscriptionFromRepo.TopicID
	subscription.UpdatedAt = time.Now().UTC()
	err = s.subscriptionRepo.Store(tx, subscription)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		rollback(tx, "subscription-update-rollback")
		return err
	}

	return nil
}

// NewSubscription returns a new Subscription with SubscriptionRepo
func NewSubscription(topicRepo hammer.TopicRepository, subscriptionRepo hammer.SubscriptionRepository, txFactoryRepo hammer.TxFactoryRepository) Subscription {
	return Subscription{
		topicRepo:        topicRepo,
		subscriptionRepo: subscriptionRepo,
		txFactoryRepo:    txFactoryRepo,
	}
}
