package service

import (
	"context"

	"github.com/allisson/hammer"
)

// Delivery is a implementation of hammer.DeliveryService
type Delivery struct {
	deliveryRepo        hammer.DeliveryRepository
	deliveryAttemptRepo hammer.DeliveryAttemptRepository
}

// Find returns hammer.Delivery by id
func (d Delivery) Find(ctx context.Context, id string) (*hammer.Delivery, error) {
	return d.deliveryRepo.Find(ctx, id)
}

// FindAll returns []hammer.Delivery by findOptions
func (d Delivery) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Delivery, error) {
	return d.deliveryRepo.FindAll(ctx, findOptions)
}

// NewDelivery returns a new Delivery with DeliveryRepo
func NewDelivery(deliveryRepo hammer.DeliveryRepository, deliveryAttemptRepo hammer.DeliveryAttemptRepository) *Delivery {
	return &Delivery{
		deliveryRepo:        deliveryRepo,
		deliveryAttemptRepo: deliveryAttemptRepo,
	}
}
