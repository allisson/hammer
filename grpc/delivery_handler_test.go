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

func TestDeliveryHandler(t *testing.T) {
	t.Run("Test GetDelivery", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		handler := NewDeliveryHandler(deliveryService)
		ctx := context.Background()
		delivery := hammer.Delivery{
			ID:        "id",
			TopicID:   "topic_id",
			Data:      "{}",
			CreatedAt: time.Now().UTC(),
		}
		request := &pb.GetDeliveryRequest{
			Id: "id",
		}
		deliveryService.On("Find", mock.Anything).Return(delivery, nil)

		response, err := handler.GetDelivery(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "id", response.Id)
		assert.Equal(t, "{}", response.Data)
	})

	t.Run("Test ListDeliveries", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		handler := NewDeliveryHandler(deliveryService)
		ctx := context.Background()
		delivery := hammer.Delivery{
			ID:        "id",
			TopicID:   "topic_id",
			Data:      "{}",
			CreatedAt: time.Now().UTC(),
		}
		request := &pb.ListDeliveriesRequest{
			Limit:  50,
			Offset: 0,
		}
		deliveryService.On("FindAll", mock.Anything).Return([]hammer.Delivery{delivery}, nil)

		response, err := handler.ListDeliveries(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(response.Deliveries))
		assert.Equal(t, "id", response.Deliveries[0].Id)
		assert.Equal(t, "{}", response.Deliveries[0].Data)
	})
}
