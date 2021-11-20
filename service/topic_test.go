package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
)

func TestTopic(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Find", func(t *testing.T) {
		expectedTopic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(expectedTopic, nil)

		topic, err := topicService.Find(ctx, expectedTopic.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedTopic, topic)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedTopics := []*hammer.Topic{hammer.MakeTestTopic()}
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedTopics, nil)

		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		topics, err := topicService.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedTopics, topics)
	})

	t.Run("Test Create", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(&hammer.Topic{}, sql.ErrNoRows)

		err := topicService.Create(ctx, topic)
		assert.Nil(t, err)
	})

	t.Run("Test Create with topic already exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(&hammer.Topic{}, nil)

		err := topicService.Create(ctx, topic)
		assert.Equal(t, hammer.ErrTopicAlreadyExists, err)
	})

	t.Run("Test Update", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(&hammer.Topic{}, nil)

		topic.Name = "My Topic"
		err := topicService.Update(ctx, topic)
		assert.Nil(t, err)
	})

	t.Run("Test Update with topic does not exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(topic, sql.ErrNoRows)

		topic.Name = "My Topic"
		err := topicService.Update(ctx, topic)
		assert.Equal(t, hammer.ErrTopicDoesNotExists, err)
	})

	t.Run("Test Delete", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Find", mock.Anything, mock.Anything).Return(topic, nil)
		topicRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

		err := topicService.Delete(ctx, topic.ID)
		assert.Nil(t, err)
	})
}
