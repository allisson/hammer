package grpc

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
)

// SubscriptionHandler implements methods for topic create/update
type SubscriptionHandler struct {
	subscriptionService hammer.SubscriptionService
}

func (s SubscriptionHandler) buildResponse(subscription *hammer.Subscription) *pb.Subscription {
	response := &pb.Subscription{}
	response.Id = subscription.ID
	response.TopicId = subscription.TopicID
	response.Name = subscription.Name
	response.Url = subscription.URL
	response.SecretToken = subscription.SecretToken
	response.MaxDeliveryAttempts = uint32(subscription.MaxDeliveryAttempts)
	response.DeliveryAttemptDelay = uint32(subscription.DeliveryAttemptDelay)
	response.DeliveryAttemptTimeout = uint32(subscription.DeliveryAttemptTimeout)
	response.CreatedAt = timestamppb.New(subscription.CreatedAt)
	response.UpdatedAt = timestamppb.New(subscription.UpdatedAt)

	return response
}

// CreateSubscription creates a new subscription
func (s SubscriptionHandler) CreateSubscription(ctx context.Context, request *pb.CreateSubscriptionRequest) (*pb.Subscription, error) {
	if request.Subscription == nil {
		request.Subscription = &pb.Subscription{}
	}

	// Build a subscription
	subscription := &hammer.Subscription{
		ID:                     request.Subscription.Id,
		TopicID:                request.Subscription.TopicId,
		Name:                   request.Subscription.Name,
		URL:                    request.Subscription.Url,
		SecretToken:            request.Subscription.SecretToken,
		MaxDeliveryAttempts:    int(request.Subscription.MaxDeliveryAttempts),
		DeliveryAttemptDelay:   int(request.Subscription.DeliveryAttemptDelay),
		DeliveryAttemptTimeout: int(request.Subscription.DeliveryAttemptTimeout),
	}

	// Validate subscription
	err := subscription.Validate()
	if err != nil {
		st := validationStatusError(codes.InvalidArgument, "invalid_subscription", err)
		return &pb.Subscription{}, st.Err()
	}

	// Create subscription
	err = s.subscriptionService.Create(ctx, subscription)
	if err != nil {
		return &pb.Subscription{}, status.Error(codes.Internal, err.Error())
	}
	response := s.buildResponse(subscription)
	return response, nil
}

// UpdateSubscription update the subscription
func (s SubscriptionHandler) UpdateSubscription(ctx context.Context, request *pb.UpdateSubscriptionRequest) (*pb.Subscription, error) {
	if request.Subscription == nil {
		request.Subscription = &pb.Subscription{}
	}

	// Build a subscription
	subscription := &hammer.Subscription{
		ID:                     request.Subscription.Id,
		TopicID:                request.Subscription.TopicId,
		Name:                   request.Subscription.Name,
		URL:                    request.Subscription.Url,
		SecretToken:            request.Subscription.SecretToken,
		MaxDeliveryAttempts:    int(request.Subscription.MaxDeliveryAttempts),
		DeliveryAttemptDelay:   int(request.Subscription.DeliveryAttemptDelay),
		DeliveryAttemptTimeout: int(request.Subscription.DeliveryAttemptTimeout),
	}

	// Validate subscription
	err := subscription.Validate()
	if err != nil {
		return &pb.Subscription{}, status.Error(codes.InvalidArgument, "invalid_subscription")
	}

	// Update subscription
	err = s.subscriptionService.Update(ctx, subscription)
	if err != nil {
		return &pb.Subscription{}, status.Error(codes.Internal, err.Error())
	}
	response := s.buildResponse(subscription)
	return response, nil
}

// GetSubscription gets the subscription
func (s SubscriptionHandler) GetSubscription(ctx context.Context, request *pb.GetSubscriptionRequest) (*pb.Subscription, error) {
	// Get subscription from service
	subscription, err := s.subscriptionService.Find(ctx, request.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &pb.Subscription{}, status.Error(codes.NotFound, hammer.ErrSubscriptionDoesNotExists.Error())
		default:
			return &pb.Subscription{}, status.Error(codes.Internal, err.Error())
		}
	}
	response := s.buildResponse(subscription)
	return response, nil
}

// ListSubscriptions get a list of topics
func (s SubscriptionHandler) ListSubscriptions(ctx context.Context, request *pb.ListSubscriptionsRequest) (*pb.ListSubscriptionsResponse, error) {
	// Get limit and offset
	limit, offset := parsePagination(request.Limit, request.Offset)

	// Create response
	response := &pb.ListSubscriptionsResponse{}

	// Get subscriptions
	findOptions := hammer.FindOptions{
		FindPagination: &hammer.FindPagination{
			Limit:  limit,
			Offset: offset,
		},
	}
	createdAtFilters := createdAtFilters(request.CreatedAtGt, request.CreatedAtGte, request.CreatedAtLt, request.CreatedAtLte)
	findOptions.FindFilters = append(findOptions.FindFilters, createdAtFilters...)
	subscriptions, err := s.subscriptionService.FindAll(ctx, findOptions)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	// Update response
	for i := range subscriptions {
		subscription := subscriptions[i]
		subscriptionResponse := s.buildResponse(subscription)
		response.Subscriptions = append(response.Subscriptions, subscriptionResponse)
	}

	return response, nil
}

// DeleteSubscription delete the subscription
func (s SubscriptionHandler) DeleteSubscription(ctx context.Context, request *pb.DeleteSubscriptionRequest) (*empty.Empty, error) {
	response := &empty.Empty{}

	// Delete topic
	err := s.subscriptionService.Delete(ctx, request.Id)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// NewSubscriptionHandler returns a new Topic
func NewSubscriptionHandler(subscriptionService hammer.SubscriptionService) SubscriptionHandler {
	return SubscriptionHandler{subscriptionService: subscriptionService}
}
