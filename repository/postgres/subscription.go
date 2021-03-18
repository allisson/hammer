package repository

import (
	"context"
	"database/sql"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// Subscription is a implementation of hammer.SubscriptionRepository
type Subscription struct {
	db *sqlx.DB
}

// Find returns hammer.Subscription by id
func (s *Subscription) Find(ctx context.Context, id string) (*hammer.Subscription, error) {
	subscription := &hammer.Subscription{}
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	query, args := findQuery("subscriptions", findOptions)
	err := s.db.GetContext(ctx, subscription, query, args...)
	return subscription, err
}

// FindAll returns []hammer.Subscription by limit and offset
func (s *Subscription) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Subscription, error) {
	subscriptions := []*hammer.Subscription{}
	query, args := findQuery("subscriptions", findOptions)
	err := s.db.SelectContext(ctx, &subscriptions, query, args...)
	return subscriptions, err
}

// Store a hammer.Subscription on database (create or update)
func (s *Subscription) Store(ctx context.Context, subscription *hammer.Subscription) error {
	_, err := s.Find(ctx, subscription.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			query, args := insertQuery("subscriptions", subscription)
			_, err := s.db.ExecContext(ctx, query, args...)
			return err
		}
		return err
	}
	query, args := updateQuery("subscriptions", subscription.ID, subscription)
	_, err = s.db.ExecContext(ctx, query, args...)
	return err
}

// Delete a hammer.Subscription on database
func (s *Subscription) Delete(ctx context.Context, id string) error {
	_, err := s.Find(ctx, id)
	if err != nil {
		return err
	}
	query, args := deleteQuery("subscriptions", id)
	_, err = s.db.ExecContext(ctx, query, args...)
	return err
}

// NewSubscription returns a new Subscription with db connection
func NewSubscription(db *sqlx.DB) *Subscription {
	return &Subscription{db: db}
}
