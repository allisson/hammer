package service

import (
	"time"

	"github.com/allisson/hammer"
)

// Subscription is a implementation of hammer.SubscriptionService
type Subscription struct {
	subscriptionRepo hammer.SubscriptionRepository
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
	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	return s.subscriptionRepo.Store(subscription)
}

// Update a hammer.Subscription on repository
func (s *Subscription) Update(subscription *hammer.Subscription) error {
	subscription.UpdatedAt = time.Now().UTC()
	return s.subscriptionRepo.Store(subscription)
}

// NewSubscription returns a new Subscription with SubscriptionRepo
func NewSubscription(subscriptionRepo hammer.SubscriptionRepository) Subscription {
	return Subscription{subscriptionRepo: subscriptionRepo}
}
