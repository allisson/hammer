package service

import (
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
		topicService := NewTopic(topicRepo)
		topicRepo.On("Find", mock.Anything).Return(expectedTopic, nil)

		topic, err := topicService.Find(expectedTopic.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedTopic, topic)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		expectedTopics := []hammer.Topic{hammer.MakeTestTopic()}
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("FindAll", mock.Anything, mock.Anything).Return(expectedTopics, nil)

		topics, err := topicService.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedTopics, topics)
	})

	t.Run("Test Create", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Store", mock.Anything).Return(nil)

		err := topicService.Create(&topic)
		assert.Nil(t, err)
	})

	t.Run("Test Update", func(t *testing.T) {
		topic := hammer.MakeTestTopic()
		topicRepo := &mocks.TopicRepository{}
		topicService := NewTopic(topicRepo)
		topicRepo.On("Store", mock.Anything).Return(nil)

		err := topicService.Create(&topic)
		assert.Nil(t, err)

		topic.Name = "My Topic"
		err = topicService.Update(&topic)
		assert.Nil(t, err)
	})
}
