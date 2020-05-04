package service

import (
	"github.com/allisson/hammer"
)

// Delivery is a implementation of hammer.DeliveryService
type Delivery struct {
	deliveryRepo hammer.DeliveryRepository
}

// Find returns hammer.Delivery by id
func (d *Delivery) Find(id string) (hammer.Delivery, error) {
	return d.deliveryRepo.Find(id)
}

// FindAll returns []hammer.Delivery by limit and offset
func (d *Delivery) FindAll(limit, offset int) ([]hammer.Delivery, error) {
	return d.deliveryRepo.FindAll(limit, offset)
}

// FindToDispatch returns []hammer.Delivery ready to dispatch by limit and offset
func (d *Delivery) FindToDispatch(limit, offset int) ([]hammer.Delivery, error) {
	return d.deliveryRepo.FindToDispatch(limit, offset)
}

// NewDelivery returns a new Delivery with DeliveryRepo
func NewDelivery(deliveryRepo hammer.DeliveryRepository) Delivery {
	return Delivery{deliveryRepo: deliveryRepo}
}
