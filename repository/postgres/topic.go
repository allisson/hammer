package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/allisson/hammer"
)

// Topic is a implementation of hammer.TopicRepository
type Topic struct {
	db *sqlx.DB
}

// Find returns hammer.Topic by id
func (t *Topic) Find(ctx context.Context, id string) (*hammer.Topic, error) {
	topic := &hammer.Topic{}
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	query, args := findQuery("topics", findOptions)
	err := t.db.GetContext(ctx, topic, query, args...)
	return topic, err
}

// FindAll returns []hammer.Topic by limit and offset
func (t *Topic) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Topic, error) {
	topics := []*hammer.Topic{}
	query, args := findQuery("topics", findOptions)
	err := t.db.SelectContext(ctx, &topics, query, args...)
	return topics, err
}

// Store a hammer.Topic on database (create or update)
func (t *Topic) Store(ctx context.Context, topic *hammer.Topic) error {
	_, err := t.Find(ctx, topic.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			query, args := insertQuery("topics", topic)
			_, err := t.db.ExecContext(ctx, query, args...)
			return err
		}
		return err
	}
	query, args := updateQuery("topics", topic.ID, topic)
	_, err = t.db.ExecContext(ctx, query, args...)
	return err
}

// Delete a hammer.Topic on database
func (t *Topic) Delete(ctx context.Context, id string) error {
	_, err := t.Find(ctx, id)
	if err != nil {
		return err
	}
	query, args := deleteQuery("topics", id)
	_, err = t.db.ExecContext(ctx, query, args...)
	return err
}

// NewTopic returns a new Topic with db connection
func NewTopic(db *sqlx.DB) *Topic {
	return &Topic{db: db}
}
