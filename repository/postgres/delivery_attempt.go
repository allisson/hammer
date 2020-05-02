package repository

import (
	"database/sql"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// DeliveryAttempt is a implementation of hammer.DeliveryAttemptRepository
type DeliveryAttempt struct {
	db *sqlx.DB
}

// Find returns hammer.DeliveryAttempt by id
func (d *DeliveryAttempt) Find(id string) (hammer.DeliveryAttempt, error) {
	deliveryAttempt := hammer.DeliveryAttempt{}
	err := d.db.Get(&deliveryAttempt, sqlDeliveryAttemptFind, id)
	return deliveryAttempt, err
}

// FindAll returns []hammer.DeliveryAttempt by limit and offset
func (d *DeliveryAttempt) FindAll(limit, offset int) ([]hammer.DeliveryAttempt, error) {
	deliveryAttempts := []hammer.DeliveryAttempt{}
	err := d.db.Select(&deliveryAttempts, sqlDeliveryAttemptFindAll, limit, offset)
	return deliveryAttempts, err
}

func (d *DeliveryAttempt) create(deliveryAttempt *hammer.DeliveryAttempt) error {
	_, err := d.db.NamedExec(sqlDeliveryAttemptCreate, deliveryAttempt)
	return err
}

func (d *DeliveryAttempt) update(deliveryAttempt *hammer.DeliveryAttempt) error {
	_, err := d.db.NamedExec(sqlDeliveryAttemptUpdate, deliveryAttempt)
	return err
}

// Store a hammer.DeliveryAttempt on database (create or update)
func (d *DeliveryAttempt) Store(deliveryAttempt *hammer.DeliveryAttempt) error {
	_, err := d.Find(deliveryAttempt.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return d.create(deliveryAttempt)
		}
		return err
	}
	return d.update(deliveryAttempt)
}

// NewDeliveryAttempt returns a new DeliveryAttempt with db connection
func NewDeliveryAttempt(db *sqlx.DB) DeliveryAttempt {
	return DeliveryAttempt{db: db}
}
