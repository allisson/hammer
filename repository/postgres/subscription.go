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
	err := s.db.Get(&subscription, sqlSubscriptionFind, id)
	return subscription, err
}

// FindAll returns []hammer.Subscription by limit and offset
func (s *Subscription) FindAll(limit, offset int) ([]hammer.Subscription, error) {
	subscriptions := []hammer.Subscription{}
	err := s.db.Select(&subscriptions, sqlSubscriptionFindAll, limit, offset)
	return subscriptions, err
}

// FindByTopic returns hammer.Subscription by topic_id and topic_created_at
func (s *Subscription) FindByTopic(topicID string) ([]hammer.Subscription, error) {
	subscriptions := []hammer.Subscription{}
	sqlStatement := strings.ReplaceAll(sqlSubscriptionFind, "WHERE id = $1", "WHERE topic_id = $1")
	err := s.db.Select(&subscriptions, sqlStatement, topicID)
	return subscriptions, err
}

func (s *Subscription) create(subscription *hammer.Subscription) error {
	_, err := s.db.NamedExec(sqlSubscriptionCreate, subscription)
	return err
}

func (s *Subscription) update(subscription *hammer.Subscription) error {
	_, err := s.db.NamedExec(sqlSubscriptionUpdate, subscription)
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
