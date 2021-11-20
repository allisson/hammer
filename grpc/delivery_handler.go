package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
)

// DeliveryHandler implements methods for Delivery get/list
type DeliveryHandler struct {
	deliveryService hammer.DeliveryService
}

func (d DeliveryHandler) buildResponse(delivery *hammer.Delivery) *pb.Delivery {
	response := &pb.Delivery{}
	response.Id = delivery.ID
	response.TopicId = delivery.TopicID
	response.SubscriptionId = delivery.SubscriptionID
	response.MessageId = delivery.MessageID
	response.ContentType = delivery.ContentType
	response.Data = delivery.Data
	response.Url = delivery.URL
	response.SecretToken = delivery.SecretToken
	response.MaxDeliveryAttempts = uint32(delivery.MaxDeliveryAttempts)
	response.DeliveryAttemptDelay = uint32(delivery.DeliveryAttemptDelay)
	response.DeliveryAttemptTimeout = uint32(delivery.DeliveryAttemptTimeout)
	response.ScheduledAt = timestamppb.New(delivery.ScheduledAt)
	response.DeliveryAttempts = uint32(delivery.DeliveryAttempts)
	response.Status = delivery.Status
	response.CreatedAt = timestamppb.New(delivery.CreatedAt)
	response.UpdatedAt = timestamppb.New(delivery.UpdatedAt)

	return response
}

// GetDelivery gets the Delivery
func (d DeliveryHandler) GetDelivery(ctx context.Context, request *pb.GetDeliveryRequest) (*pb.Delivery, error) {
	// Get delivery from service
	delivery, err := d.deliveryService.Find(ctx, request.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &pb.Delivery{}, status.Error(codes.NotFound, hammer.ErrDeliveryDoesNotExists.Error())
		default:
			return &pb.Delivery{}, status.Error(codes.Internal, err.Error())
		}
	}
	response := d.buildResponse(delivery)
	return response, nil
}

// ListDeliveries get a list of deliveries
func (d DeliveryHandler) ListDeliveries(ctx context.Context, request *pb.ListDeliveriesRequest) (*pb.ListDeliveriesResponse, error) {
	// Get limit and offset
	limit, offset := parsePagination(request.Limit, request.Offset)

	// Create response
	response := &pb.ListDeliveriesResponse{}

	// Get Delivery
	findOptions := hammer.FindOptions{
		FindPagination: &hammer.FindPagination{
			Limit:  limit,
			Offset: offset,
		},
	}
	if request.TopicId != "" {
		topicFilter := hammer.FindFilter{
			FieldName: "topic_id",
			Operator:  "=",
			Value:     request.TopicId,
		}
		findOptions.FindFilters = append(findOptions.FindFilters, topicFilter)
	}
	if request.SubscriptionId != "" {
		subscriptionFilter := hammer.FindFilter{
			FieldName: "subscription_id",
			Operator:  "=",
			Value:     request.SubscriptionId,
		}
		findOptions.FindFilters = append(findOptions.FindFilters, subscriptionFilter)
	}
	if request.MessageId != "" {
		messageFilter := hammer.FindFilter{
			FieldName: "message_id",
			Operator:  "=",
			Value:     request.MessageId,
		}
		findOptions.FindFilters = append(findOptions.FindFilters, messageFilter)
	}
	if request.Status != "" {
		statusFilter := hammer.FindFilter{
			FieldName: "status",
			Operator:  "=",
			Value:     request.Status,
		}
		findOptions.FindFilters = append(findOptions.FindFilters, statusFilter)
	}
	createdAtFilters := createdAtFilters(request.CreatedAtGt, request.CreatedAtGte, request.CreatedAtLt, request.CreatedAtLte)
	findOptions.FindFilters = append(findOptions.FindFilters, createdAtFilters...)
	deliveries, err := d.deliveryService.FindAll(ctx, findOptions)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	// Update response
	for i := range deliveries {
		delivery := deliveries[i]
		deliveryResponse := d.buildResponse(delivery)
		response.Deliveries = append(response.Deliveries, deliveryResponse)
	}

	return response, nil
}

// NewDeliveryHandler returns a new Delivery
func NewDeliveryHandler(deliveryService hammer.DeliveryService) DeliveryHandler {
	return DeliveryHandler{deliveryService: deliveryService}
}
