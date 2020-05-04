package service

import (
	"database/sql"
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopic(t *testing.T) {
	t.Run("Test Find", func(t *testing.T) {
		expectedTopic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		topicService := NewTopic(topicRepo, txFactoryRepo)
		topicRepo.On("Find", mock.Anything).Return(expectedTopic, nil)

		topic, err := topicService.Find(expectedTopic.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedTopic, topic)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedTopics := []hammer.Topic{hammer.MakeTestTopic()}
		topicRepo := &mocks.TopicRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		topicService := NewTopic(topicRepo, txFactoryRepo)
		topicRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedTopics, nil)

		topics, err := topicService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedTopics, topics)
	})

	t.Run("Test Create", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		topicService := NewTopic(topicRepo, txFactoryRepo)
		txFactoryRepo.On("New").Return(txRepo, nil)
		topicRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, sql.ErrNoRows)
		txRepo.On("Commit").Return(nil)

		err := topicService.Create(&topic)
		assert.Nil(t, err)
	})

	t.Run("Test Create with topic already exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		topicService := NewTopic(topicRepo, txFactoryRepo)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, nil)

		err := topicService.Create(&topic)
		assert.Equal(t, hammer.ErrTopicAlreadyExists, err)
	})

	t.Run("Test Update", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		txRepo := &mocks.TxRepository{}
		topicService := NewTopic(topicRepo, txFactoryRepo)
		txFactoryRepo.On("New").Return(txRepo, nil)
		topicRepo.On("Store", mock.Anything, mock.Anything).Return(nil)
		topicRepo.On("Find", mock.Anything).Return(hammer.Topic{}, nil)
		txRepo.On("Commit").Return(nil)

		topic.Name = "My Topic"
		err := topicService.Update(&topic)
		assert.Nil(t, err)
	})

	t.Run("Test Update with topic does not exists on repository", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		txFactoryRepo := &mocks.TxFactoryRepository{}
		topicService := NewTopic(topicRepo, txFactoryRepo)
		topicRepo.On("Find", mock.Anything).Return(topic, sql.ErrNoRows)

		topic.Name = "My Topic"
		err := topicService.Update(&topic)
		assert.Equal(t, hammer.ErrTopicDoesNotExists, err)
	})
}
