package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/allisson/hammer"
)

// DeliveryAttempt is a implementation of hammer.DeliveryAttemptRepository
type DeliveryAttempt struct {
	db *sqlx.DB
}

// Find returns hammer.DeliveryAttempt by id
func (d DeliveryAttempt) Find(ctx context.Context, id string) (*hammer.DeliveryAttempt, error) {
	deliveryAttempt := &hammer.DeliveryAttempt{}
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	query, args := findQuery("delivery_attempts", findOptions)
	err := d.db.GetContext(ctx, deliveryAttempt, query, args...)
	return deliveryAttempt, err
}

// FindAll returns []hammer.DeliveryAttempt by limit and offset
func (d DeliveryAttempt) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.DeliveryAttempt, error) {
	deliveryAttempts := []*hammer.DeliveryAttempt{}
	query, args := findQuery("delivery_attempts", findOptions)
	err := d.db.SelectContext(ctx, &deliveryAttempts, query, args...)
	return deliveryAttempts, err
}

// Store a hammer.DeliveryAttempt on database (create or update)
func (d DeliveryAttempt) Store(ctx context.Context, deliveryAttempt *hammer.DeliveryAttempt) error {
	_, err := d.Find(ctx, deliveryAttempt.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			query, args := insertQuery("delivery_attempts", deliveryAttempt)
			_, err := d.db.ExecContext(ctx, query, args...)
			return err
		}
		return err
	}
	query, args := updateQuery("delivery_attempts", deliveryAttempt.ID, deliveryAttempt)
	_, err = d.db.ExecContext(ctx, query, args...)
	return err
}

// NewDeliveryAttempt returns a new DeliveryAttempt with db connection
func NewDeliveryAttempt(db *sqlx.DB) *DeliveryAttempt {
	return &DeliveryAttempt{db: db}
}
