package grpc

import (
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/allisson/hammer/mocks"
)

func TestTopicHandler(t *testing.T) {
	t.Run("Test CreateTopic", func(t *testing.T) {
		topicService := &mocks.TopicService{}
		handler := NewTopicHandler(topicService)
		ctx := context.Background()
		request := &pb.CreateTopicRequest{
			Topic: &pb.Topic{
				Id:   "topic_id",
				Name: "Topic",
			},
		}
		topicService.On("Create", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.CreateTopic(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "topic_id", response.Id)
		assert.Equal(t, "Topic", response.Name)
	})

	t.Run("Test UpdateTopic", func(t *testing.T) {
		topicService := &mocks.TopicService{}
		handler := NewTopicHandler(topicService)
		ctx := context.Background()
		request := &pb.UpdateTopicRequest{
			Topic: &pb.Topic{
				Id:   "topic_id",
				Name: "Topic",
			},
		}
		topicService.On("Update", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.UpdateTopic(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "topic_id", response.Id)
		assert.Equal(t, "Topic", response.Name)
	})

	t.Run("Test GetTopic", func(t *testing.T) {
		topicService := &mocks.TopicService{}
		handler := NewTopicHandler(topicService)
		ctx := context.Background()
		topic := &hammer.Topic{
			ID:   "topic_id",
			Name: "Topic",
		}
		request := &pb.GetTopicRequest{
			Id: "topic_id",
		}
		topicService.On("Find", mock.Anything, mock.Anything).Return(topic, nil)

		response, err := handler.GetTopic(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "topic_id", response.Id)
		assert.Equal(t, "Topic", response.Name)
	})

	t.Run("Test ListTopics", func(t *testing.T) {
		topicService := &mocks.TopicService{}
		handler := NewTopicHandler(topicService)
		ctx := context.Background()
		topic := &hammer.Topic{
			ID:   "topic_id",
			Name: "Topic",
		}
		request := &pb.ListTopicsRequest{
			Limit:  50,
			Offset: 0,
		}
		topicService.On("FindAll", mock.Anything, mock.Anything).Return([]*hammer.Topic{topic}, nil)

		response, err := handler.ListTopics(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(response.Topics))
		assert.Equal(t, "topic_id", response.Topics[0].Id)
		assert.Equal(t, "Topic", response.Topics[0].Name)
	})

	t.Run("Test Delete", func(t *testing.T) {
		topicService := &mocks.TopicService{}
		handler := NewTopicHandler(topicService)
		ctx := context.Background()
		topic := hammer.Topic{
			ID:   "topic_id",
			Name: "Topic",
		}
		request := &pb.DeleteTopicRequest{
			Id: topic.ID,
		}
		topicService.On("Delete", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.DeleteTopic(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, &empty.Empty{}, response)
	})
}
