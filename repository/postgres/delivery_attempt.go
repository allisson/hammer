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
	data := map[string]interface{}{
		"table": "delivery_attempts",
		"id":    id,
	}
	query, args, err := buildQuery(sqlFind, data)
	if err != nil {
		return deliveryAttempt, err
	}
	err = d.db.Get(&deliveryAttempt, query, args...)
	return deliveryAttempt, err
}

// FindAll returns []hammer.DeliveryAttempt by limit and offset
func (d *DeliveryAttempt) FindAll(limit, offset int) ([]hammer.DeliveryAttempt, error) {
	deliveryAttempts := []hammer.DeliveryAttempt{}
	data := map[string]interface{}{
		"table":   "delivery_attempts",
		"limit":   limit,
		"offset":  offset,
		"orderBy": "id DESC",
	}
	query, args, err := buildQuery(sqlFindAll, data)
	if err != nil {
		return deliveryAttempts, err
	}
	err = d.db.Select(&deliveryAttempts, query, args...)
	return deliveryAttempts, err
}

// Store a hammer.DeliveryAttempt on database (create or update)
func (d *DeliveryAttempt) Store(tx hammer.TxRepository, deliveryAttempt *hammer.DeliveryAttempt) error {
	_, err := d.Find(deliveryAttempt.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tx.Exec(sqlDeliveryAttemptCreate, deliveryAttempt)
		}
		return err
	}
	return tx.Exec(sqlDeliveryAttemptUpdate, deliveryAttempt)
}

// NewDeliveryAttempt returns a new DeliveryAttempt with db connection
func NewDeliveryAttempt(db *sqlx.DB) DeliveryAttempt {
	return DeliveryAttempt{db: db}
}
