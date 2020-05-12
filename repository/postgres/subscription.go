package repository

import (
	"database/sql"
	"strings"

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
	data := map[string]interface{}{
		"table": "subscriptions",
		"id":    id,
	}
	query, args, err := buildQuery(sqlFind, data)
	if err != nil {
		return subscription, err
	}
	err = s.db.Get(&subscription, query, args...)
	return subscription, err
}

// FindAll returns []hammer.Subscription by limit and offset
func (s *Subscription) FindAll(limit, offset int) ([]hammer.Subscription, error) {
	subscriptions := []hammer.Subscription{}
	data := map[string]interface{}{
		"table":   "subscriptions",
		"limit":   limit,
		"offset":  offset,
		"orderBy": "id ASC",
	}
	query, args, err := buildQuery(sqlFindAll, data)
	if err != nil {
		return subscriptions, err
	}
	err = s.db.Select(&subscriptions, query, args...)
	return subscriptions, err
}

// FindByTopic returns hammer.Subscription by topic_id and topic_created_at
func (s *Subscription) FindByTopic(topicID string) ([]hammer.Subscription, error) {
	subscriptions := []hammer.Subscription{}
	sqlStatement := strings.ReplaceAll(sqlSubscriptionFind, "WHERE id = $1", "WHERE topic_id = $1")
	err := s.db.Select(&subscriptions, sqlStatement, topicID)
	return subscriptions, err
}

// Store a hammer.Subscription on database (create or update)
func (s *Subscription) Store(tx hammer.TxRepository, subscription *hammer.Subscription) error {
	_, err := s.Find(subscription.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tx.Exec(sqlSubscriptionCreate, subscription)
		}
		return err
	}
	return tx.Exec(sqlSubscriptionUpdate, subscription)
}

// NewSubscription returns a new Subscription with db connection
func NewSubscription(db *sqlx.DB) Subscription {
	return Subscription{db: db}
}
