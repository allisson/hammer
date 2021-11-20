package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/allisson/hammer"
)

// Message is a implementation of hammer.MessageRepository
type Message struct {
	db *sqlx.DB
}

func (m Message) create(ctx context.Context, message *hammer.Message) error {
	// Start transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	// Create message
	query, args := insertQuery("messages", message)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollback(tx)
		return err
	}

	// Get subscriptions
	subscriptions := []*hammer.Subscription{}
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "topic_id",
				Operator:  "=",
				Value:     message.TopicID,
			},
		},
	}
	query, args = findQuery("subscriptions", findOptions)
	if err := tx.SelectContext(ctx, &subscriptions, query, args...); err != nil {
		rollback(tx)
		return err
	}

	// Create deliveries
	for _, subscription := range subscriptions {
		id, err := hammer.GenerateULID()
		if err != nil {
			rollback(tx)
			return err
		}
		now := time.Now().UTC()
		delivery := hammer.Delivery{
			ID:                     id,
			TopicID:                message.TopicID,
			SubscriptionID:         subscription.ID,
			MessageID:              message.ID,
			ContentType:            message.ContentType,
			Data:                   message.Data,
			URL:                    subscription.URL,
			SecretToken:            subscription.SecretToken,
			MaxDeliveryAttempts:    subscription.MaxDeliveryAttempts,
			DeliveryAttemptDelay:   subscription.DeliveryAttemptDelay,
			DeliveryAttemptTimeout: subscription.DeliveryAttemptTimeout,
			ScheduledAt:            now,
			Status:                 hammer.DeliveryStatusPending,
			CreatedAt:              now,
			UpdatedAt:              now,
		}
		query, args := insertQuery("deliveries", delivery)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			rollback(tx)
			return err
		}
	}

	return tx.Commit()
}

// Find returns hammer.Message by id
func (m Message) Find(ctx context.Context, id string) (*hammer.Message, error) {
	message := &hammer.Message{}
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	query, args := findQuery("messages", findOptions)
	err := m.db.GetContext(ctx, message, query, args...)
	return message, err
}

// FindAll returns []hammer.Message by limit and offset
func (m Message) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Message, error) {
	messages := []*hammer.Message{}
	query, args := findQuery("messages", findOptions)
	err := m.db.SelectContext(ctx, &messages, query, args...)
	return messages, err
}

// Store a hammer.Message on database (create or update)
func (m Message) Store(ctx context.Context, message *hammer.Message) error {
	_, err := m.Find(ctx, message.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return m.create(ctx, message)
		}
		return err
	}
	query, args := updateQuery("messages", message.ID, message)
	_, err = m.db.ExecContext(ctx, query, args...)
	return err
}

// Delete a hammer.Message on database
func (m Message) Delete(ctx context.Context, id string) error {
	_, err := m.Find(ctx, id)
	if err != nil {
		return err
	}
	query, args := deleteQuery("messages", id)
	_, err = m.db.ExecContext(ctx, query, args...)
	return err
}

// NewMessage returns a new Message with db connection
func NewMessage(db *sqlx.DB) *Message {
	return &Message{db: db}
}
