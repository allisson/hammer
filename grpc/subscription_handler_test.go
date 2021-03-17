package grpc

import (
	"testing"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/allisson/hammer/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestSubscriptionHandler(t *testing.T) {
	t.Run("Test CreateSubscription", func(t *testing.T) {
		subscriptionService := &mocks.SubscriptionService{}
		handler := NewSubscriptionHandler(subscriptionService)
		ctx := context.Background()
		request := &pb.CreateSubscriptionRequest{
			Subscription: &pb.Subscription{
				Id:                     "subscription_id",
				TopicId:                "topic_id",
				Name:                   "Subscription",
				Url:                    "https://example.com/post",
				MaxDeliveryAttempts:    5,
				DeliveryAttemptDelay:   60,
				DeliveryAttemptTimeout: 5,
			},
		}
		subscriptionService.On("Create", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.CreateSubscription(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "subscription_id", response.Id)
		assert.Equal(t, "Subscription", response.Name)
	})

	t.Run("Test UpdateSubscription", func(t *testing.T) {
		subscriptionService := &mocks.SubscriptionService{}
		handler := NewSubscriptionHandler(subscriptionService)
		ctx := context.Background()
		request := &pb.UpdateSubscriptionRequest{
			Subscription: &pb.Subscription{
				Id:                     "subscription_id",
				TopicId:                "topic_id",
				Name:                   "Subscription",
				Url:                    "https://example.com/post",
				MaxDeliveryAttempts:    5,
				DeliveryAttemptDelay:   60,
				DeliveryAttemptTimeout: 5,
			},
		}
		subscriptionService.On("Update", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.UpdateSubscription(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "subscription_id", response.Id)
		assert.Equal(t, "Subscription", response.Name)
	})

	t.Run("Test GetSubscription", func(t *testing.T) {
		subscriptionService := &mocks.SubscriptionService{}
		handler := NewSubscriptionHandler(subscriptionService)
		ctx := context.Background()
		subscription := &hammer.Subscription{
			ID:   "subscription_id",
			Name: "Subscription",
		}
		request := &pb.GetSubscriptionRequest{
			Id: "subscription_id",
		}
		subscriptionService.On("Find", mock.Anything, mock.Anything).Return(subscription, nil)

		response, err := handler.GetSubscription(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, "subscription_id", response.Id)
		assert.Equal(t, "Subscription", response.Name)
	})

	t.Run("Test ListSubscriptions", func(t *testing.T) {
		subscriptionService := &mocks.SubscriptionService{}
		handler := NewSubscriptionHandler(subscriptionService)
		ctx := context.Background()
		subscription := &hammer.Subscription{
			ID:   "subscription_id",
			Name: "Subscription",
		}
		request := &pb.ListSubscriptionsRequest{
			Limit:  50,
			Offset: 0,
		}
		subscriptionService.On("FindAll", mock.Anything, mock.Anything).Return([]*hammer.Subscription{subscription}, nil)

		response, err := handler.ListSubscriptions(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(response.Subscriptions))
		assert.Equal(t, "subscription_id", response.Subscriptions[0].Id)
		assert.Equal(t, "Subscription", response.Subscriptions[0].Name)
	})

	t.Run("Test Delete", func(t *testing.T) {
		subscriptionService := &mocks.SubscriptionService{}
		handler := NewSubscriptionHandler(subscriptionService)
		ctx := context.Background()
		subscription := hammer.Subscription{
			ID:   "subscription_id",
			Name: "Subscription",
		}
		request := &pb.DeleteSubscriptionRequest{
			Id: subscription.ID,
		}
		subscriptionService.On("Delete", mock.Anything, mock.Anything).Return(nil)

		response, err := handler.DeleteSubscription(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, &empty.Empty{}, response)
	})
}
