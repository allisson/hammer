package grpc

import (
	"context"
	"database/sql"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SubscriptionHandler implements methods for topic create/update
type SubscriptionHandler struct {
	subscriptionService hammer.SubscriptionService
}

func (s *SubscriptionHandler) buildResponse(subscription *hammer.Subscription) (*pb.Subscription, error) {
	response := &pb.Subscription{}
	createdAt, err := ptypes.TimestampProto(subscription.CreatedAt)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}
	updatedAt, err := ptypes.TimestampProto(subscription.UpdatedAt)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}
	response.Id = subscription.ID
	response.TopicId = subscription.TopicID
	response.Name = subscription.Name
	response.Url = subscription.URL
	response.SecretToken = subscription.SecretToken
	response.MaxDeliveryAttempts = uint32(subscription.MaxDeliveryAttempts)
	response.DeliveryAttemptDelay = uint32(subscription.DeliveryAttemptDelay)
	response.DeliveryAttemptTimeout = uint32(subscription.DeliveryAttemptTimeout)
	response.CreatedAt = createdAt
	response.UpdatedAt = updatedAt

	return response, nil
}

// CreateSubscription creates a new subscription
func (s *SubscriptionHandler) CreateSubscription(ctx context.Context, request *pb.CreateSubscriptionRequest) (*pb.Subscription, error) {
	if request.Subscription == nil {
		request.Subscription = &pb.Subscription{}
	}

	// Build a subscription
	subscription := hammer.Subscription{
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

	// Create subscription
	err = s.subscriptionService.Create(&subscription)
	if err != nil {
		return &pb.Subscription{}, status.Error(codes.Internal, err.Error())
	}

	return s.buildResponse(&subscription)
}

// GetSubscription gets the subscription
func (s *SubscriptionHandler) GetSubscription(ctx context.Context, request *pb.GetSubscriptionRequest) (*pb.Subscription, error) {
	// Get subscription from service
	subscription, err := s.subscriptionService.Find(request.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &pb.Subscription{}, status.Error(codes.NotFound, hammer.ErrSubscriptionDoesNotExists.Error())
		default:
			return &pb.Subscription{}, status.Error(codes.Internal, err.Error())
		}
	}

	return s.buildResponse(&subscription)
}

// ListSubscriptions get a list of topics
func (s *SubscriptionHandler) ListSubscriptions(ctx context.Context, request *pb.ListSubscriptionsRequest) (*pb.ListSubscriptionsResponse, error) {
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
	subscriptions, err := s.subscriptionService.FindAll(findOptions)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	// Update response
	for _, subscription := range subscriptions {
		subscriptionResponse, err := s.buildResponse(&subscription)
		if err != nil {
			return response, status.Error(codes.Internal, err.Error())
		}
		response.Subscriptions = append(response.Subscriptions, subscriptionResponse)
	}

	return response, nil
}

// NewSubscriptionHandler returns a new Topic
func NewSubscriptionHandler(subscriptionService hammer.SubscriptionService) SubscriptionHandler {
	return SubscriptionHandler{subscriptionService: subscriptionService}
}
