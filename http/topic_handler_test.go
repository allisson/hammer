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

func TestTopicHandler(t *testing.T) {
	t.Run("Test Post with malformed request body", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Post)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"malformed_request_body","details":"unexpected end of JSON input"}`).
			End()
	})

	t.Run("Test Post with invalid request body", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Post)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"invalid_request_body","details":"{\"id\":\"cannot be blank\",\"name\":\"cannot be blank\"}"}`).
			End()
	})

	t.Run("Test Post with object already exists error", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Post)
		topicService.On("Create", mock.Anything).Return(hammer.ErrObjectAlreadyExists)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{"id": "topic", "name": "Topic"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"object_already_exists","details":""}`).
			End()
	})

	t.Run("Test Post with service error", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Post)
		topicService.On("Create", mock.Anything).Return(errors.New("service_error"))

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{"id": "topic", "name": "Topic"}`).
			Expect(t).
			Status(http.StatusInternalServerError).
			Body(`{"message":"service_error","details":"service_error"}`).
			End()
	})

	t.Run("Test Post success", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Post)
		topicService.On("Create", mock.Anything).Return(nil)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{"id": "topic", "name": "Topic"}`).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})
}
