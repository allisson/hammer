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

// TopicHandler implements methods for topic create/update
type TopicHandler struct {
	topicService hammer.TopicService
}

func (t *TopicHandler) buildResponse(topic *hammer.Topic) (*pb.Topic, error) {
	response := &pb.Topic{}
	createdAt, err := ptypes.TimestampProto(topic.CreatedAt)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}
	updatedAt, err := ptypes.TimestampProto(topic.UpdatedAt)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}
	response.Id = topic.ID
	response.Name = topic.Name
	response.CreatedAt = createdAt
	response.UpdatedAt = updatedAt

	return response, nil
}

// CreateTopic creates a new topic
func (t *TopicHandler) CreateTopic(ctx context.Context, request *pb.CreateTopicRequest) (*pb.Topic, error) {
	if request.Topic == nil {
		request.Topic = &pb.Topic{}
	}

	// Build a topic
	topic := hammer.Topic{
		ID:   request.Topic.Id,
		Name: request.Topic.Name,
	}

	// Validate topic
	err := topic.Validate()
	if err != nil {
		return &pb.Topic{}, status.Error(codes.InvalidArgument, "invalid_topic")
	}

	// Create topic
	err = t.topicService.Create(&topic)
	if err != nil {
		return &pb.Topic{}, status.Error(codes.Internal, err.Error())
	}

	return t.buildResponse(&topic)
}

// GetTopic gets the topic
func (t *TopicHandler) GetTopic(ctx context.Context, request *pb.GetTopicRequest) (*pb.Topic, error) {
	// Get topic from service
	topic, err := t.topicService.Find(request.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &pb.Topic{}, status.Error(codes.NotFound, hammer.ErrTopicDoesNotExists.Error())
		default:
			return &pb.Topic{}, status.Error(codes.Internal, err.Error())
		}
	}

	return t.buildResponse(&topic)
}

// ListTopics get a list of topics
func (t *TopicHandler) ListTopics(ctx context.Context, request *pb.ListTopicsRequest) (*pb.ListTopicsResponse, error) {
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
	topics, err := t.topicService.FindAll(findOptions)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	// Update response
	for _, topic := range topics {
		topicResponse, err := t.buildResponse(&topic)
		if err != nil {
			return response, status.Error(codes.Internal, err.Error())
		}
		response.Topics = append(response.Topics, topicResponse)
	}

	return response, nil
}

// NewTopicHandler returns a new Topic
func NewTopicHandler(topicService hammer.TopicService) TopicHandler {
	return TopicHandler{topicService: topicService}
}
