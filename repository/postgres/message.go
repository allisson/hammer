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
	err := m.db.Get(&message, sqlMessageFind, id)
	return message, err
}

// FindAll returns []hammer.Message by limit and offset
func (m *Message) FindAll(limit, offset int) ([]hammer.Message, error) {
	messages := []hammer.Message{}
	err := m.db.Select(&messages, sqlMessageFindAll, limit, offset)
	return messages, err
}

// FindByTopic returns []hammer.Message by topic, limit and offset
func (m *Message) FindByTopic(topicID string, limit, offset int) ([]hammer.Message, error) {
	messages := []hammer.Message{}
	err := m.db.Select(&messages, sqlMessageFindByTopic, topicID, limit, offset)
	return messages, err
}

func (m *Message) create(message *hammer.Message) error {
	_, err := m.db.NamedExec(sqlMessageCreate, message)
	return err
}

func (m *Message) update(message *hammer.Message) error {
	_, err := m.db.NamedExec(sqlMessageUpdate, message)
	return err
}

// Store a hammer.Message on database (create or update)
func (m *Message) Store(message *hammer.Message) error {
	_, err := m.Find(message.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return m.create(message)
		}
		return err
	}
	return m.update(message)
}

// NewMessage returns a new Message with db connection
func NewMessage(db *sqlx.DB) Message {
	return Message{db: db}
}
