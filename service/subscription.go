package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/allisson/hammer"
)

// Subscription is a implementation of hammer.SubscriptionService
type Subscription struct {
	topicRepo        hammer.TopicRepository
	subscriptionRepo hammer.SubscriptionRepository
}

func (s Subscription) topicExists(ctx context.Context, topicID string) error {
	_, err := s.topicRepo.Find(ctx, topicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrTopicDoesNotExists
		}
		return err
	}
	return nil
}

// Find returns hammer.Subscription by id
func (s Subscription) Find(ctx context.Context, id string) (*hammer.Subscription, error) {
	return s.subscriptionRepo.Find(ctx, id)
}

// FindAll returns []hammer.Subscription by findOptions
func (s Subscription) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Subscription, error) {
	return s.subscriptionRepo.FindAll(ctx, findOptions)
}

// Create a hammer.Subscription on repository
func (s Subscription) Create(ctx context.Context, subscription *hammer.Subscription) error {
	// Verify if subscription already exists
	_, err := s.subscriptionRepo.Find(ctx, subscription.ID)
	if err == nil {
		return hammer.ErrSubscriptionAlreadyExists
	}

	// Verify if topic already exists
	err = s.topicExists(ctx, subscription.TopicID)
	if err != nil {
		return err
	}

	// Create new subscription with default values
	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	if subscription.SecretToken == "" {
		subscription.SecretToken = generateRandomString(hammer.DefaultSecretTokenLength)
	}
	err = s.subscriptionRepo.Store(ctx, subscription)
	if err != nil {
		return err
	}

	return nil
}

// Update a hammer.Subscription on repository
func (s Subscription) Update(ctx context.Context, subscription *hammer.Subscription) error {
	// Verify if subscription already exists
	subscriptionFromRepo, err := s.subscriptionRepo.Find(ctx, subscription.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrSubscriptionDoesNotExists
		}
		return err
	}

	// Verify if topic already exists
	err = s.topicExists(ctx, subscription.TopicID)
	if err != nil {
		return err
	}

	// Update subscription
	subscription.ID = subscriptionFromRepo.ID
	subscription.TopicID = subscriptionFromRepo.TopicID
	subscription.CreatedAt = subscriptionFromRepo.CreatedAt
	subscription.UpdatedAt = time.Now().UTC()
	err = s.subscriptionRepo.Store(ctx, subscription)
	if err != nil {
		return err
	}

	return nil
}

// Delete a hammer.Subscription on repository
func (s Subscription) Delete(ctx context.Context, id string) error {
	return s.subscriptionRepo.Delete(ctx, id)
}

// NewSubscription returns a new Subscription with SubscriptionRepo
func NewSubscription(topicRepo hammer.TopicRepository, subscriptionRepo hammer.SubscriptionRepository) *Subscription {
	return &Subscription{
		topicRepo:        topicRepo,
		subscriptionRepo: subscriptionRepo,
	}
}
