package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/allisson/hammer"
)

// TopicHandler implements methods for topic create/update
type TopicHandler struct {
	topicService hammer.TopicService
}

// Create a new topic
func (t *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"

	// Parse request
	topic := hammer.Topic{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "read_body_error", Details: err.Error()}, contentType)
		return
	}
	err = json.Unmarshal(requestBody, &topic)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: "malformed_request_body", Details: err.Error()}, contentType)
		return
	}

	// Validate topic
	err = topic.Validate()
	if err != nil {
		errorPayload, _ := json.Marshal(err)
		errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: "invalid_request_body", Details: string(errorPayload)}, contentType)
		return
	}

	// Call service
	err = t.topicService.Create(&topic)
	if err != nil {
		switch err {
		case hammer.ErrTopicAlreadyExists:
			errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: err.Error(), Details: ""}, contentType)
			return
		default:
			errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "service_error", Details: err.Error()}, contentType)
			return
		}
	}

	// Convert to json
	responseBody, err := json.Marshal(&topic)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "json_convert_error", Details: err.Error()}, contentType)
		return
	}

	makeResponse(w, responseBody, http.StatusCreated, contentType)
}

// List topics
func (t *TopicHandler) List(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"

	// Get limit and offset
	limit, offset := getLimitOffset(r)

	// Call service
	topics, err := t.topicService.FindAll(limit, offset)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "service_error", Details: err.Error()}, contentType)
		return
	}

	// Create ListTopicsResponse
	topicResponse := hammer.ListTopicsResponse{
		Limit:  limit,
		Offset: offset,
		Topics: topics,
	}

	// Convert to json
	responseBody, err := json.Marshal(&topicResponse)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "json_convert_error", Details: err.Error()}, contentType)
		return
	}

	makeResponse(w, responseBody, http.StatusOK, contentType)
}

// NewTopicHandler returns a new TopicHandler
func NewTopicHandler(topicService hammer.TopicService) TopicHandler {
	return TopicHandler{topicService: topicService}
}
