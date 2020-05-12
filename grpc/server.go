package grpc

import (
	"context"

	pb "github.com/allisson/hammer/api/v1"
)

// Server implements grpc server
type Server struct {
	topicHandler           TopicHandler
	subscriptionHandler    SubscriptionHandler
	messageHandler         MessageHandler
	deliveryHandler        DeliveryHandler
	deliveryAttemptHandler DeliveryAttemptHandler
}

// CreateTopic creates a new topic
func (s *Server) CreateTopic(ctx context.Context, request *pb.CreateTopicRequest) (*pb.Topic, error) {
	return s.topicHandler.CreateTopic(ctx, request)
}

// GetTopic gets the topic
func (s *Server) GetTopic(ctx context.Context, request *pb.GetTopicRequest) (*pb.Topic, error) {
	return s.topicHandler.GetTopic(ctx, request)
}

// ListTopics get a list of topics
func (s *Server) ListTopics(ctx context.Context, request *pb.ListTopicsRequest) (*pb.ListTopicsResponse, error) {
	return s.topicHandler.ListTopics(ctx, request)
}

// CreateSubscription creates a new subscription
func (s *Server) CreateSubscription(ctx context.Context, request *pb.CreateSubscriptionRequest) (*pb.Subscription, error) {
	return s.subscriptionHandler.CreateSubscription(ctx, request)
}

// GetSubscription gets the subscription
func (s *Server) GetSubscription(ctx context.Context, request *pb.GetSubscriptionRequest) (*pb.Subscription, error) {
	return s.subscriptionHandler.GetSubscription(ctx, request)
}

// ListSubscriptions get a list of subscriptions
func (s *Server) ListSubscriptions(ctx context.Context, request *pb.ListSubscriptionsRequest) (*pb.ListSubscriptionsResponse, error) {
	return s.subscriptionHandler.ListSubscriptions(ctx, request)
}

// CreateMessage creates a new message
func (s *Server) CreateMessage(ctx context.Context, request *pb.CreateMessageRequest) (*pb.Message, error) {
	return s.messageHandler.CreateMessage(ctx, request)
}

// GetMessage gets the message
func (s *Server) GetMessage(ctx context.Context, request *pb.GetMessageRequest) (*pb.Message, error) {
	return s.messageHandler.GetMessage(ctx, request)
}

// ListMessages get a list of messages
func (s *Server) ListMessages(ctx context.Context, request *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	return s.messageHandler.ListMessages(ctx, request)
}

// GetDelivery gets the delivery
func (s *Server) GetDelivery(ctx context.Context, request *pb.GetDeliveryRequest) (*pb.Delivery, error) {
	return s.deliveryHandler.GetDelivery(ctx, request)
}

// ListDeliveries get a list of deliveries
func (s *Server) ListDeliveries(ctx context.Context, request *pb.ListDeliveriesRequest) (*pb.ListDeliveriesResponse, error) {
	return s.deliveryHandler.ListDeliveries(ctx, request)
}

// GetDeliveryAttempt gets the delivery attempt
func (s *Server) GetDeliveryAttempt(ctx context.Context, request *pb.GetDeliveryAttemptRequest) (*pb.DeliveryAttempt, error) {
	return s.deliveryAttemptHandler.GetDeliveryAttempt(ctx, request)
}

// ListDeliveryAttempts get a list of delivery attempts
func (s *Server) ListDeliveryAttempts(ctx context.Context, request *pb.ListDeliveryAttemptsRequest) (*pb.ListDeliveryAttemptsResponse, error) {
	return s.deliveryAttemptHandler.ListDeliveryAttempts(ctx, request)
}

// NewServer returns a new server
func NewServer(topicHandler TopicHandler, subscriptionHandler SubscriptionHandler, messageHandler MessageHandler, deliveryHandler DeliveryHandler, deliveryAttemptHandler DeliveryAttemptHandler) Server {
	return Server{
		topicHandler:           topicHandler,
		subscriptionHandler:    subscriptionHandler,
		messageHandler:         messageHandler,
		deliveryHandler:        deliveryHandler,
		deliveryAttemptHandler: deliveryAttemptHandler,
	}
}
