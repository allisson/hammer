package grpc

import (
	"testing"
	"time"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestDeliveryAttemptHandler(t *testing.T) {
	t.Run("Test GetDeliveryAttempt", func(t *testing.T) {
		deliveryAttemptService := &mocks.DeliveryAttemptService{}
		handler := NewDeliveryAttemptHandler(deliveryAttemptService)
		ctx := context.Background()
		deliveryAttempt := &hammer.DeliveryAttempt{
			ID:         "id",
			DeliveryID: "delivery_id",
			CreatedAt:  time.Now().UTC(),
		}
		request := &pb.GetDeliveryAttemptRequest{
			Id: "id",
		}
		deliveryAttemptService.On("Find", mock.Anything, mock.Anything).Return(deliveryAttempt, nil)

		response, err := handler.GetDeliveryAttempt(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "id", response.Id)
		assert.Equal(t, "delivery_id", response.DeliveryId)
	})

	t.Run("Test ListDeliveryAttempts", func(t *testing.T) {
		deliveryAttemptService := &mocks.DeliveryAttemptService{}
		handler := NewDeliveryAttemptHandler(deliveryAttemptService)
		ctx := context.Background()
		deliveryAttempt := &hammer.DeliveryAttempt{
			ID:         "id",
			DeliveryID: "delivery_id",
			CreatedAt:  time.Now().UTC(),
		}
		request := &pb.ListDeliveryAttemptsRequest{
			Limit:  50,
			Offset: 0,
		}
		deliveryAttemptService.On("FindAll", mock.Anything, mock.Anything).Return([]*hammer.DeliveryAttempt{deliveryAttempt}, nil)

		response, err := handler.ListDeliveryAttempts(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(response.DeliveryAttempts))
		assert.Equal(t, "id", response.DeliveryAttempts[0].Id)
		assert.Equal(t, "delivery_id", response.DeliveryAttempts[0].DeliveryId)
	})
}
