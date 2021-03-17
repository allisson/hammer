package service

import (
	"context"

	"github.com/allisson/hammer"
)

// DeliveryAttempt is a implementation of hammer.DeliveryAttemptService
type DeliveryAttempt struct {
	deliveryAttemptRepo hammer.DeliveryAttemptRepository
}

// Find returns hammer.DeliveryAttempt by id
func (d DeliveryAttempt) Find(ctx context.Context, id string) (*hammer.DeliveryAttempt, error) {
	return d.deliveryAttemptRepo.Find(ctx, id)
}

// FindAll returns []hammer.DeliveryAttempt by findOptions
func (d DeliveryAttempt) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.DeliveryAttempt, error) {
	return d.deliveryAttemptRepo.FindAll(ctx, findOptions)
}

// NewDeliveryAttempt returns a new DeliveryAttempt with DeliveryAttemptRepo
func NewDeliveryAttempt(deliveryAttemptRepo hammer.DeliveryAttemptRepository) *DeliveryAttempt {
	return &DeliveryAttempt{deliveryAttemptRepo: deliveryAttemptRepo}
}
