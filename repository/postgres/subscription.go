package postgres

import (
	"database/sql"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// Subscription is a implementation of hammer.SubscriptionRepository
type Subscription struct {
	db *sqlx.DB
}

// Find returns hammer.Subscription by id
func (s *Subscription) Find(id string) (hammer.Subscription, error) {
	subscription := hammer.Subscription{}
	sqlStatement := `
		SELECT *
		FROM subscriptions
		WHERE id = $1
	`
	err := s.db.Get(&subscription, sqlStatement, id)
	return subscription, err
}

func (s *Subscription) create(subscription *hammer.Subscription) error {
	sqlStatement := `
		INSERT INTO subscriptions (
			"id",
			"topic_id",
			"name",
			"url",
			"secret_token",
			"max_delivery_attempts",
			"delivery_attempt_delay",
			"delivery_attempt_timeout",
			"active",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:topic_id,
			:name,
			:url,
			:secret_token,
			:max_delivery_attempts,
			:delivery_attempt_delay,
			:delivery_attempt_timeout,
			:active,
			:created_at,
			:updated_at
		)
	`
	_, err := s.db.NamedExec(sqlStatement, subscription)
	return err
}

func (s *Subscription) update(subscription *hammer.Subscription) error {
	sqlStatement := `
		UPDATE subscriptions
		SET topic_id = :topic_id,
			name = :name,
			url = :url,
			secret_token = :secret_token,
			max_delivery_attempts = :max_delivery_attempts,
			delivery_attempt_delay = :delivery_attempt_delay,
			delivery_attempt_timeout = :delivery_attempt_timeout,
			active = :active,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := s.db.NamedExec(sqlStatement, subscription)
	return err
}

// Store a hammer.Subscription on database (create or update)
func (s *Subscription) Store(subscription *hammer.Subscription) error {
	_, err := s.Find(subscription.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.create(subscription)
		}
		return err
	}
	return s.update(subscription)
}

// NewSubscription returns a new Subscription with db connection
func NewSubscription(db *sqlx.DB) Subscription {
	return Subscription{db: db}
}
