package grpc

import (
	"context"
	"database/sql"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MessageHandler implements methods for message create/update
type MessageHandler struct {
	messageService hammer.MessageService
}

func (m *MessageHandler) buildResponse(message *hammer.Message) (*pb.Message, error) {
	response := &pb.Message{}
	createdAt, err := ptypes.TimestampProto(message.CreatedAt)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}
	response.Id = message.ID
	response.TopicId = message.TopicID
	response.ContentType = message.ContentType
	response.Data = message.Data
	response.CreatedAt = createdAt

	return response, nil
}

// CreateMessage creates a new Message
func (m *MessageHandler) CreateMessage(ctx context.Context, request *pb.CreateMessageRequest) (*pb.Message, error) {
	if request.Message == nil {
		request.Message = &pb.Message{}
	}

	// Build a message
	message := hammer.Message{
		ID:          request.Message.Id,
		TopicID:     request.Message.TopicId,
		ContentType: request.Message.ContentType,
		Data:        string(request.Message.Data),
	}

	// Validate message
	err := message.Validate()
	if err != nil {
		st := validationStatusError(codes.InvalidArgument, "invalid_message", err)
		return &pb.Message{}, st.Err()
	}

	// Create Message
	err = m.messageService.Create(&message)
	if err != nil {
		return &pb.Message{}, status.Error(codes.Internal, err.Error())
	}

	return m.buildResponse(&message)
}

// GetMessage gets the message
func (m *MessageHandler) GetMessage(ctx context.Context, request *pb.GetMessageRequest) (*pb.Message, error) {
	// Get nessage from service
	message, err := m.messageService.Find(request.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &pb.Message{}, status.Error(codes.NotFound, hammer.ErrMessageDoesNotExists.Error())
		default:
			return &pb.Message{}, status.Error(codes.Internal, err.Error())
		}
	}

	return m.buildResponse(&message)
}

// ListMessages get a list of messages
func (m *MessageHandler) ListMessages(ctx context.Context, request *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	// Get limit and offset
	limit, offset := parsePagination(request.Limit, request.Offset)

	// Create response
	response := &pb.ListMessagesResponse{}

	// Get messages
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
	createdAtFilters := createdAtFilters(request.CreatedAtGt, request.CreatedAtGte, request.CreatedAtLt, request.CreatedAtLte)
	findOptions.FindFilters = append(findOptions.FindFilters, createdAtFilters...)
	messages, err := m.messageService.FindAll(findOptions)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	// Update response
	for _, message := range messages {
		messageResponse, err := m.buildResponse(&message)
		if err != nil {
			return response, status.Error(codes.Internal, err.Error())
		}
		response.Messages = append(response.Messages, messageResponse)
	}

	return response, nil
}

// DeleteMessage delete the message
func (m *MessageHandler) DeleteMessage(ctx context.Context, request *pb.DeleteMessageRequest) (*empty.Empty, error) {
	response := &empty.Empty{}

	// Delete topic
	err := m.messageService.Delete(request.Id)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// NewMessageHandler returns a new Message
func NewMessageHandler(messageService hammer.MessageService) MessageHandler {
	return MessageHandler{messageService: messageService}
}
