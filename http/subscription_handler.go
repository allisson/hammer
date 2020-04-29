package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/allisson/hammer"
)

// SubscriptionHandler implements methods for Subscription create/update
type SubscriptionHandler struct {
	subscriptionService hammer.SubscriptionService
}

// Create new subscription
func (s *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"

	// Parse request
	subscription := hammer.Subscription{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "read_body_error", Details: err.Error()}, contentType)
		return
	}
	err = json.Unmarshal(requestBody, &subscription)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: "malformed_request_body", Details: err.Error()}, contentType)
		return
	}

	// Validate Subscription
	err = subscription.Validate()
	if err != nil {
		errorPayload, _ := json.Marshal(err)
		errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: "invalid_request_body", Details: string(errorPayload)}, contentType)
		return
	}

	// Call service
	err = s.subscriptionService.Create(&subscription)
	if err != nil {
		switch err {
		case hammer.ErrTopicDoesNotExists, hammer.ErrSubscriptionAlreadyExists:
			errorResponse(w, hammer.Error{Code: http.StatusBadRequest, Message: err.Error(), Details: ""}, contentType)
			return
		default:
			errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "service_error", Details: err.Error()}, contentType)
			return
		}
	}

	// Convert to json
	responseBody, err := json.Marshal(&subscription)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "json_convert_error", Details: err.Error()}, contentType)
		return
	}

	makeResponse(w, responseBody, http.StatusCreated, contentType)
}

// List subscriptions
func (s *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"

	// Get limit and offset
	limit, offset := getLimitOffset(r)

	// Call service
	subscriptions, err := s.subscriptionService.FindAll(limit, offset)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "service_error", Details: err.Error()}, contentType)
		return
	}

	// Create ListSubscriptionsResponse
	subscriptionResponse := hammer.ListSubscriptionsResponse{
		Limit:         limit,
		Offset:        offset,
		Subscriptions: subscriptions,
	}

	// Convert to json
	responseBody, err := json.Marshal(&subscriptionResponse)
	if err != nil {
		errorResponse(w, hammer.Error{Code: http.StatusInternalServerError, Message: "json_convert_error", Details: err.Error()}, contentType)
		return
	}

	makeResponse(w, responseBody, http.StatusOK, contentType)
}

// NewSubscriptionHandler returns a new SubscriptionHandler
func NewSubscriptionHandler(subscriptionService hammer.SubscriptionService) SubscriptionHandler {
	return SubscriptionHandler{subscriptionService: subscriptionService}
}
