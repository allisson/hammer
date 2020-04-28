package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/allisson/hammer"
	"github.com/allisson/hammer/mocks"
	"github.com/go-chi/chi"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionHandler(t *testing.T) {
	t.Run("Test Post with malformed request body", func(t *testing.T) {
		subscriptionService := mocks.SubscriptionService{}
		subscriptionHandler := NewSubscriptionHandler(&subscriptionService)
		r := chi.NewRouter()
		r.Post("/subscriptions", subscriptionHandler.Post)

		apitest.New().
			Handler(r).
			Post("/subscriptions").
			Body(`{`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"malformed_request_body","details":"unexpected end of JSON input"}`).
			End()
	})

	t.Run("Test Post with invalid request body", func(t *testing.T) {
		subscriptionService := mocks.SubscriptionService{}
		subscriptionHandler := NewSubscriptionHandler(&subscriptionService)
		r := chi.NewRouter()
		r.Post("/subscriptions", subscriptionHandler.Post)

		apitest.New().
			Handler(r).
			Post("/subscriptions").
			Body(`{}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"invalid_request_body","details":"{\"id\":\"cannot be blank\",\"name\":\"cannot be blank\",\"topic_id\":\"cannot be blank\",\"url\":\"cannot be blank\"}"}`).
			End()
	})

	t.Run("Test Post with subscription already exists error", func(t *testing.T) {
		subscriptionService := mocks.SubscriptionService{}
		subscriptionHandler := NewSubscriptionHandler(&subscriptionService)
		r := chi.NewRouter()
		r.Post("/subscriptions", subscriptionHandler.Post)
		subscriptionService.On("Create", mock.Anything).Return(hammer.ErrSubscriptionAlreadyExists)

		apitest.New().
			Handler(r).
			Post("/subscriptions").
			Body(`{"id": "subscription", "topic_id": "topic", "name": "subscription", "url": "http://example.com"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"subscription_already_exists","details":""}`).
			End()
	})

	t.Run("Test Post with topic does not exists error", func(t *testing.T) {
		subscriptionService := mocks.SubscriptionService{}
		subscriptionHandler := NewSubscriptionHandler(&subscriptionService)
		r := chi.NewRouter()
		r.Post("/subscriptions", subscriptionHandler.Post)
		subscriptionService.On("Create", mock.Anything).Return(hammer.ErrTopicDoesNotExists)

		apitest.New().
			Handler(r).
			Post("/subscriptions").
			Body(`{"id": "subscription", "topic_id": "topic", "name": "subscription", "url": "http://example.com"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"topic_does_not_exists","details":""}`).
			End()
	})

	t.Run("Test Post with service error", func(t *testing.T) {
		subscriptionService := mocks.SubscriptionService{}
		subscriptionHandler := NewSubscriptionHandler(&subscriptionService)
		r := chi.NewRouter()
		r.Post("/subscriptions", subscriptionHandler.Post)
		subscriptionService.On("Create", mock.Anything).Return(errors.New("service_error"))

		apitest.New().
			Handler(r).
			Post("/subscriptions").
			Body(`{"id": "subscription", "topic_id": "topic", "name": "subscription", "url": "http://example.com"}`).
			Expect(t).
			Status(http.StatusInternalServerError).
			Body(`{"message":"service_error","details":"service_error"}`).
			End()
	})

	t.Run("Test Post success", func(t *testing.T) {
		subscriptionService := mocks.SubscriptionService{}
		subscriptionHandler := NewSubscriptionHandler(&subscriptionService)
		r := chi.NewRouter()
		r.Post("/subscriptions", subscriptionHandler.Post)
		subscriptionService.On("Create", mock.Anything).Return(nil)

		apitest.New().
			Handler(r).
			Post("/subscriptions").
			Body(`{"id": "subscription", "topic_id": "topic", "name": "subscription", "url": "http://example.com"}`).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})
}
