package repository

import (
	"database/sql"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// Topic is a implementation of hammer.TopicRepository
type Topic struct {
	db *sqlx.DB
}

// Find returns hammer.Topic by id
func (t *Topic) Find(id string) (hammer.Topic, error) {
	topic := hammer.Topic{}
	sqlStatement := `
		SELECT *
		FROM topics
		WHERE id = $1
	`
	err := t.db.Get(&topic, sqlStatement, id)
	return topic, err
}

// FindAll returns []hammer.Topic by limit and offset
func (t *Topic) FindAll(limit, offset int) ([]hammer.Topic, error) {
	topics := []hammer.Topic{}
	sqlStatement := `
		SELECT *
		FROM topics
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
	`
	err := t.db.Select(&topics, sqlStatement, limit, offset)
	return topics, err
}

func (t *Topic) create(topic *hammer.Topic) error {
	sqlStatement := `
		INSERT INTO topics (
			"id",
			"name",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:name,
			:created_at,
			:updated_at
		)
	`
	_, err := t.db.NamedExec(sqlStatement, topic)
	return err
}

func (t *Topic) update(topic *hammer.Topic) error {
	sqlStatement := `
		UPDATE topics
		SET name = :name,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := t.db.NamedExec(sqlStatement, topic)
	return err
}

// Store a hammer.Topic on database (create or update)
func (t *Topic) Store(topic *hammer.Topic) error {
	_, err := t.Find(topic.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return t.create(topic)
		}
		return err
	}
	return t.update(topic)
}

// NewTopic returns a new Topic with db connection
func NewTopic(db *sqlx.DB) Topic {
	return Topic{db: db}
}
