package grpc

import (
	"testing"
	"time"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/allisson/hammer/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestMessageHandler(t *testing.T) {
	t.Run("Test CreateMessage", func(t *testing.T) {
		messageService := &mocks.MessageService{}
		handler := NewMessageHandler(messageService)
		ctx := context.Background()
		request := &pb.CreateMessageRequest{
			Message: &pb.Message{
				Id:          "id",
				TopicId:     "topic_id",
				ContentType: "application/json",
				Data:        "{}",
			},
		}
		messageService.On("Create", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.CreateMessage(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "id", response.Id)
		assert.Equal(t, "{}", response.Data)
	})

	t.Run("Test GetMessage", func(t *testing.T) {
		messageService := &mocks.MessageService{}
		handler := NewMessageHandler(messageService)
		ctx := context.Background()
		message := &hammer.Message{
			ID:          "id",
			TopicID:     "topic_id",
			ContentType: "application/json",
			Data:        "{}",
			CreatedAt:   time.Now().UTC(),
		}
		request := &pb.GetMessageRequest{
			Id: "id",
		}
		messageService.On("Find", mock.Anything, mock.Anything).Return(message, nil)

		response, err := handler.GetMessage(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "id", response.Id)
		assert.Equal(t, "{}", response.Data)
	})

	t.Run("Test ListMessages", func(t *testing.T) {
		messageService := &mocks.MessageService{}
		handler := NewMessageHandler(messageService)
		ctx := context.Background()
		message := &hammer.Message{
			ID:          "id",
			TopicID:     "topic_id",
			ContentType: "application/json",
			Data:        "{}",
			CreatedAt:   time.Now().UTC(),
		}
		request := &pb.ListMessagesRequest{
			Limit:  50,
			Offset: 0,
		}
		messageService.On("FindAll", mock.Anything, mock.Anything).Return([]*hammer.Message{message}, nil)

		response, err := handler.ListMessages(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(response.Messages))
		assert.Equal(t, "id", response.Messages[0].Id)
		assert.Equal(t, "{}", response.Messages[0].Data)
	})

	t.Run("Test Delete", func(t *testing.T) {
		messageService := &mocks.MessageService{}
		handler := NewMessageHandler(messageService)
		ctx := context.Background()
		message := hammer.Message{
			ID:          "id",
			TopicID:     "topic_id",
			ContentType: "application/json",
			Data:        "{}",
			CreatedAt:   time.Now().UTC(),
		}
		request := &pb.DeleteMessageRequest{
			Id: message.ID,
		}
		messageService.On("Delete", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.DeleteMessage(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, &empty.Empty{}, response)
	})
}
