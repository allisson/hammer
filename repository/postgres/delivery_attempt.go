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
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	sql, args := buildSQLQuery("delivery_attempts", findOptions)
	err := d.db.Get(&deliveryAttempt, sql, args...)
	return deliveryAttempt, err
}

// FindAll returns []hammer.DeliveryAttempt by limit and offset
func (d *DeliveryAttempt) FindAll(findOptions hammer.FindOptions) ([]hammer.DeliveryAttempt, error) {
	deliveryAttempts := []hammer.DeliveryAttempt{}
	sql, args := buildSQLQuery("delivery_attempts", findOptions)
	err := d.db.Select(&deliveryAttempts, sql, args...)
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
