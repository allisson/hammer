package service

import (
	"github.com/allisson/hammer"
)

// DeliveryAttempt is a implementation of hammer.DeliveryAttemptService
type DeliveryAttempt struct {
	deliveryAttemptRepo hammer.DeliveryAttemptRepository
}

// Find returns hammer.DeliveryAttempt by id
func (d *DeliveryAttempt) Find(id string) (hammer.DeliveryAttempt, error) {
	return d.deliveryAttemptRepo.Find(id)
}

// FindAll returns []hammer.DeliveryAttempt by findOptions
func (d *DeliveryAttempt) FindAll(findOptions hammer.FindOptions) ([]hammer.DeliveryAttempt, error) {
	return d.deliveryAttemptRepo.FindAll(findOptions)
}

// NewDeliveryAttempt returns a new DeliveryAttempt with DeliveryAttemptRepo
func NewDeliveryAttempt(deliveryAttemptRepo hammer.DeliveryAttemptRepository) DeliveryAttempt {
	return DeliveryAttempt{deliveryAttemptRepo: deliveryAttemptRepo}
}
