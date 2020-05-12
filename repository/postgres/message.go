package repository

import (
	"database/sql"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// Message is a implementation of hammer.MessageRepository
type Message struct {
	db *sqlx.DB
}

// Find returns hammer.Message by id
func (m *Message) Find(id string) (hammer.Message, error) {
	message := hammer.Message{}
	data := map[string]interface{}{
		"table": "messages",
		"id":    id,
	}
	query, args, err := buildQuery(sqlFind, data)
	if err != nil {
		return message, err
	}
	err = m.db.Get(&message, query, args...)
	return message, err
}

// FindAll returns []hammer.Message by limit and offset
func (m *Message) FindAll(limit, offset int) ([]hammer.Message, error) {
	messages := []hammer.Message{}
	data := map[string]interface{}{
		"limit":   limit,
		"offset":  offset,
		"orderBy": "id DESC",
	}
	query, args, err := buildQuery(sqlMessageFindAll, data)
	if err != nil {
		return messages, err
	}
	err = m.db.Select(&messages, query, args...)
	return messages, err
}

// FindByTopic returns []hammer.Message by topic, limit and offset
func (m *Message) FindByTopic(topicID string, limit, offset int) ([]hammer.Message, error) {
	messages := []hammer.Message{}
	data := map[string]interface{}{
		"topic_id": topicID,
		"limit":    limit,
		"offset":   offset,
		"orderBy":  "id DESC",
	}
	query, args, err := buildQuery(sqlMessageFindAll, data)
	if err != nil {
		return messages, err
	}
	err = m.db.Select(&messages, query, args...)
	return messages, err
}

// Store a hammer.Message on database (create or update)
func (m *Message) Store(tx hammer.TxRepository, message *hammer.Message) error {
	_, err := m.Find(message.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tx.Exec(sqlMessageCreate, message)
		}
		return err
	}
	return tx.Exec(sqlMessageUpdate, message)
}

// NewMessage returns a new Message with db connection
func NewMessage(db *sqlx.DB) Message {
	return Message{db: db}
}
