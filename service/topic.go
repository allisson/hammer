package service

import (
	"database/sql"
	"time"

	"github.com/allisson/hammer"
)

// Topic is a implementation of hammer.TopicService
type Topic struct {
	topicRepo hammer.TopicRepository
}

// Find returns hammer.Topic by id
func (t *Topic) Find(id string) (hammer.Topic, error) {
	return t.topicRepo.Find(id)
}

// FindAll returns []hammer.Topic by limit and offset
func (t *Topic) FindAll(limit, offset int) ([]hammer.Topic, error) {
	return t.topicRepo.FindAll(limit, offset)
}

// Create a hammer.Topic on repository
func (t *Topic) Create(topic *hammer.Topic) error {
	// Verify if object already exists
	_, err := t.topicRepo.Find(topic.ID)
	if err == nil {
		return hammer.ErrObjectAlreadyExists
	}

	// Create new topic
	now := time.Now().UTC()
	topic.CreatedAt = now
	topic.UpdatedAt = now
	return t.topicRepo.Store(topic)
}

// Update a hammer.Topic on repository
func (t *Topic) Update(topic *hammer.Topic) error {
	// Verify if object already exists
	topicFromRepo, err := t.topicRepo.Find(topic.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrObjectDoesNotExist
		}
		return err
	}

	// Update topic
	topic.ID = topicFromRepo.ID
	topic.UpdatedAt = time.Now().UTC()
	return t.topicRepo.Store(topic)
}

// NewTopic returns a new Topic with topicRepo
func NewTopic(topicRepo hammer.TopicRepository) Topic {
	return Topic{topicRepo: topicRepo}
}
