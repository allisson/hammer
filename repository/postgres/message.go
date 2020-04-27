package postgres

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
	sqlStatement := `
		SELECT *
		FROM messages
		WHERE id = $1
	`
	err := m.db.Get(&message, sqlStatement, id)
	return message, err
}

// FindAll returns hammer.Message by limit and offset
func (m *Message) FindAll(limit, offset int) ([]hammer.Message, error) {
	messages := []hammer.Message{}
	sqlStatement := `
		SELECT *
		FROM messages
		ORDER BY id DESC
		LIMIT $1
		OFFSET $2
	`
	err := m.db.Select(&messages, sqlStatement, limit, offset)
	return messages, err
}

func (m *Message) create(message *hammer.Message) error {
	sqlStatement := `
		INSERT INTO messages (
			"id",
			"topic_id",
			"data",
			"created_deliveries",
			"created_at"
		)
		VALUES (
			:id,
			:topic_id,
			:data,
			:created_deliveries,
			:created_at
		)
	`
	_, err := m.db.NamedExec(sqlStatement, message)
	return err
}

func (m *Message) update(message *hammer.Message) error {
	sqlStatement := `
		UPDATE messages
		SET topic_id = :topic_id,
			data = :data,
			created_deliveries = :created_deliveries,
			created_at = :created_at
		WHERE id = :id
	`
	_, err := m.db.NamedExec(sqlStatement, message)
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
