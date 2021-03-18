package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/allisson/hammer"
)

// Topic is a implementation of hammer.TopicService
type Topic struct {
	topicRepo hammer.TopicRepository
}

// Find returns hammer.Topic by id
func (t Topic) Find(ctx context.Context, id string) (*hammer.Topic, error) {
	return t.topicRepo.Find(ctx, id)
}

// FindAll returns []hammer.Topic by limit and offset
func (t Topic) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Topic, error) {
	return t.topicRepo.FindAll(ctx, findOptions)
}

// Create a hammer.Topic on repository
func (t Topic) Create(ctx context.Context, topic *hammer.Topic) error {
	// Verify if object already exists
	_, err := t.topicRepo.Find(ctx, topic.ID)
	if err == nil {
		return hammer.ErrTopicAlreadyExists
	}

	// Create new topic
	now := time.Now().UTC()
	topic.CreatedAt = now
	topic.UpdatedAt = now
	err = t.topicRepo.Store(ctx, topic)
	if err != nil {
		return err
	}

	return nil
}

// Update a hammer.Topic on repository
func (t Topic) Update(ctx context.Context, topic *hammer.Topic) error {
	// Verify if object already exists
	topicFromRepo, err := t.topicRepo.Find(ctx, topic.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return hammer.ErrTopicDoesNotExists
		}
		return err
	}

	// Update topic
	topic.ID = topicFromRepo.ID
	topic.CreatedAt = topicFromRepo.CreatedAt
	topic.UpdatedAt = time.Now().UTC()
	err = t.topicRepo.Store(ctx, topic)
	if err != nil {
		return err
	}

	return nil
}

// Delete a hammer.Topic on repository
func (t Topic) Delete(ctx context.Context, id string) error {
	return t.topicRepo.Delete(ctx, id)
}

// NewTopic returns a new Topic with topicRepo
func NewTopic(topicRepo hammer.TopicRepository) *Topic {
	return &Topic{
		topicRepo: topicRepo,
	}
}
