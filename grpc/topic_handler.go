package grpc

import (
	"context"
	"database/sql"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TopicHandler implements methods for topic create/update
type TopicHandler struct {
	topicService hammer.TopicService
}

func (t TopicHandler) buildResponse(topic *hammer.Topic) *pb.Topic {
	response := &pb.Topic{}
	response.Id = topic.ID
	response.Name = topic.Name
	response.CreatedAt = timestamppb.New(topic.CreatedAt)
	response.UpdatedAt = timestamppb.New(topic.UpdatedAt)

	return response
}

// CreateTopic creates a new topic
func (t TopicHandler) CreateTopic(ctx context.Context, request *pb.CreateTopicRequest) (*pb.Topic, error) {
	if request.Topic == nil {
		request.Topic = &pb.Topic{}
	}

	// Build a topic
	topic := &hammer.Topic{
		ID:   request.Topic.Id,
		Name: request.Topic.Name,
	}

	// Validate topic
	err := topic.Validate()
	if err != nil {
		st := validationStatusError(codes.InvalidArgument, "invalid_topic", err)
		return &pb.Topic{}, st.Err()
	}

	// Create topic
	err = t.topicService.Create(ctx, topic)
	if err != nil {
		return &pb.Topic{}, status.Error(codes.Internal, err.Error())
	}
	response := t.buildResponse(topic)
	return response, nil
}

// UpdateTopic update the topic
func (t TopicHandler) UpdateTopic(ctx context.Context, request *pb.UpdateTopicRequest) (*pb.Topic, error) {
	if request.Topic == nil {
		request.Topic = &pb.Topic{}
	}

	// Build a topic
	topic := &hammer.Topic{
		ID:   request.Topic.Id,
		Name: request.Topic.Name,
	}

	// Validate topic
	err := topic.Validate()
	if err != nil {
		return &pb.Topic{}, status.Error(codes.InvalidArgument, "invalid_topic")
	}

	// Update topic
	err = t.topicService.Update(ctx, topic)
	if err != nil {
		return &pb.Topic{}, status.Error(codes.Internal, err.Error())
	}
	response := t.buildResponse(topic)
	return response, nil
}

// GetTopic gets the topic
func (t TopicHandler) GetTopic(ctx context.Context, request *pb.GetTopicRequest) (*pb.Topic, error) {
	// Get topic from service
	topic, err := t.topicService.Find(ctx, request.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &pb.Topic{}, status.Error(codes.NotFound, hammer.ErrTopicDoesNotExists.Error())
		default:
			return &pb.Topic{}, status.Error(codes.Internal, err.Error())
		}
	}
	response := t.buildResponse(topic)
	return response, nil
}

// ListTopics get a list of topics
func (t TopicHandler) ListTopics(ctx context.Context, request *pb.ListTopicsRequest) (*pb.ListTopicsResponse, error) {
	// Get limit and offset
	limit, offset := parsePagination(request.Limit, request.Offset)

	// Create response
	response := &pb.ListTopicsResponse{}

	// Get topics
	findOptions := hammer.FindOptions{
		FindPagination: &hammer.FindPagination{
			Limit:  limit,
			Offset: offset,
		},
	}
	createdAtFilters := createdAtFilters(request.CreatedAtGt, request.CreatedAtGte, request.CreatedAtLt, request.CreatedAtLte)
	findOptions.FindFilters = append(findOptions.FindFilters, createdAtFilters...)
	topics, err := t.topicService.FindAll(ctx, findOptions)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	// Update response
	for i := range topics {
		topic := topics[i]
		topicResponse := t.buildResponse(topic)
		response.Topics = append(response.Topics, topicResponse)
	}

	return response, nil
}

// DeleteTopic delete the topic
func (t TopicHandler) DeleteTopic(ctx context.Context, request *pb.DeleteTopicRequest) (*empty.Empty, error) {
	response := &empty.Empty{}

	// Delete topic
	err := t.topicService.Delete(ctx, request.Id)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// NewTopicHandler returns a new Topic
func NewTopicHandler(topicService hammer.TopicService) TopicHandler {
	return TopicHandler{topicService: topicService}
}
