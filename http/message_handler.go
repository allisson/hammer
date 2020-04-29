package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/allisson/hammer"
)

// MessageHandler implements methods for Message create/update
type MessageHandler struct {
	messageService hammer.MessageService
}

// Create a new Message
func (m *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"

	// Parse request
	message := hammer.Message{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "read_body_error", Details: err.Error()}, contentType)
		return
	}
	err = json.Unmarshal(requestBody, &message)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: "malformed_request_body", Details: err.Error()}, contentType)
		return
	}

	// Validate Message
	err = message.Validate()
	if err != nil {
		errorPayload, _ := json.Marshal(err)
		errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: "invalid_request_body", Details: string(errorPayload)}, contentType)
		return
	}

	// Call service
	err = m.messageService.Create(&message)
	if err != nil {
		switch err {
		case hammer.ErrTopicDoesNotExists:
			errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: err.Error(), Details: ""}, contentType)
			return
		default:
			errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "service_error", Details: err.Error()}, contentType)
			return
		}
	}

	// Convert to json
	responseBody, err := json.Marshal(&message)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "json_convert_error", Details: err.Error()}, contentType)
		return
	}

	makeResponse(w, responseBody, http.StatusCreated, contentType)
}

// List Messages
func (m *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"

	// Get limit and offset
	limit, offset := getLimitOffset(r)

	// Call service
	messages, err := m.messageService.FindAll(limit, offset)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "service_error", Details: err.Error()}, contentType)
		return
	}

	// Create ListMessagesResponse
	MessageResponse := hammer.ListMessagesResponse{
		Limit:    limit,
		Offset:   offset,
		Messages: messages,
	}

	// Convert to json
	responseBody, err := json.Marshal(&MessageResponse)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "json_convert_error", Details: err.Error()}, contentType)
		return
	}

	makeResponse(w, responseBody, http.StatusOK, contentType)
}

// NewMessageHandler returns a new MessageHandler
func NewMessageHandler(messageService hammer.MessageService) MessageHandler {
	return MessageHandler{messageService: messageService}
}
