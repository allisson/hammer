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
	sqlStatement := `
		SELECT *
		FROM delivery_attempts
		WHERE id = $1
	`
	err := d.db.Get(&deliveryAttempt, sqlStatement, id)
	return deliveryAttempt, err
}

// FindAll returns []hammer.DeliveryAttempt by limit and offset
func (d *DeliveryAttempt) FindAll(limit, offset int) ([]hammer.DeliveryAttempt, error) {
	deliveryAttempts := []hammer.DeliveryAttempt{}
	sqlStatement := `
		SELECT *
		FROM delivery_attempts
		ORDER BY id DESC
		LIMIT $1
		OFFSET $2
	`
	err := d.db.Select(&deliveryAttempts, sqlStatement, limit, offset)
	return deliveryAttempts, err
}

func (d *DeliveryAttempt) create(deliveryAttempt *hammer.DeliveryAttempt) error {
	sqlStatement := `
		INSERT INTO delivery_attempts (
			"id",
			"delivery_id",
			"url",
			"request_headers",
			"request_body",
			"response_headers",
			"response_body",
			"response_status_code",
			"execution_duration",
			"success",
			"created_at"
		)
		VALUES (
			:id,
			:delivery_id,
			:url,
			:request_headers,
			:request_body,
			:response_headers,
			:response_body,
			:response_status_code,
			:execution_duration,
			:success,
			:created_at
		)
	`
	_, err := d.db.NamedExec(sqlStatement, deliveryAttempt)
	return err
}

func (d *DeliveryAttempt) update(deliveryAttempt *hammer.DeliveryAttempt) error {
	sqlStatement := `
		UPDATE delivery_attempts
		SET delivery_id = :delivery_id,
			url = :url,
			request_headers = :request_headers,
			request_body = :request_body,
			response_headers = :response_headers,
			response_body = :response_body,
			response_status_code = :response_status_code,
			execution_duration = :execution_duration,
			success = :success,
			created_at = :created_at
		WHERE id = :id
	`
	_, err := d.db.NamedExec(sqlStatement, deliveryAttempt)
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
