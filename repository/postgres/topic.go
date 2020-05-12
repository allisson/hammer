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
	findOptions := hammer.FindOptions{
		FindFilters: []hammer.FindFilter{
			{
				FieldName: "id",
				Operator:  "=",
				Value:     id,
			},
		},
	}
	sql, args := buildSQLQuery("topics", findOptions)
	err := t.db.Get(&topic, sql, args...)
	return topic, err
}

// FindAll returns []hammer.Topic by limit and offset
func (t *Topic) FindAll(findOptions hammer.FindOptions) ([]hammer.Topic, error) {
	topics := []hammer.Topic{}
	sql, args := buildSQLQuery("topics", findOptions)
	err := t.db.Select(&topics, sql, args...)
	return topics, err
}

// Store a hammer.Topic on database (create or update)
func (t *Topic) Store(tx hammer.TxRepository, topic *hammer.Topic) error {
	_, err := t.Find(topic.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tx.Exec(sqlTopicCreate, topic)
		}
		return err
	}
	return tx.Exec(sqlTopicUpdate, topic)
}

// NewTopic returns a new Topic with db connection
func NewTopic(db *sqlx.DB) Topic {
	return Topic{db: db}
}
