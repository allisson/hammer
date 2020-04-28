package postgres

import (
	"database/sql"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// Delivery is a implementation of hammer.DeliveryRepository
type Delivery struct {
	db *sqlx.DB
}

// Find returns hammer.Delivery by id
func (d *Delivery) Find(id string) (hammer.Delivery, error) {
	delivery := hammer.Delivery{}
	sqlStatement := `
		SELECT *
		FROM deliveries
		WHERE id = $1
	`
	err := d.db.Get(&delivery, sqlStatement, id)
	return delivery, err
}

// FindAll returns []hammer.Delivery by limit and offset
func (d *Delivery) FindAll(limit, offset int) ([]hammer.Delivery, error) {
	deliveries := []hammer.Delivery{}
	sqlStatement := `
		SELECT *
		FROM deliveries
		ORDER BY id DESC
		LIMIT $1
		OFFSET $2
	`
	err := d.db.Select(&deliveries, sqlStatement, limit, offset)
	return deliveries, err
}

func (d *Delivery) create(delivery *hammer.Delivery) error {
	sqlStatement := `
		INSERT INTO deliveries (
			"id",
			"topic_id",
			"subscription_id",
			"message_id",
			"scheduled_at",
			"delivery_attempts",
			"status",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:topic_id,
			:subscription_id,
			:message_id,
			:scheduled_at,
			:delivery_attempts,
			:status,
			:created_at,
			:updated_at
		)
	`
	_, err := d.db.NamedExec(sqlStatement, delivery)
	return err
}

func (d *Delivery) update(delivery *hammer.Delivery) error {
	sqlStatement := `
		UPDATE deliveries
		SET topic_id = :topic_id,
			subscription_id = :subscription_id,
			message_id = :message_id,
			scheduled_at = :scheduled_at,
			delivery_attempts = :delivery_attempts,
			status = :status,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := d.db.NamedExec(sqlStatement, delivery)
	return err
}

// Store a hammer.Delivery on database (create or update)
func (d *Delivery) Store(delivery *hammer.Delivery) error {
	_, err := d.Find(delivery.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return d.create(delivery)
		}
		return err
	}
	return d.update(delivery)
}

// NewDelivery returns a new Delivery with db connection
func NewDelivery(db *sqlx.DB) Delivery {
	return Delivery{db: db}
}
