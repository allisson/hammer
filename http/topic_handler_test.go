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

func TestTopicHandlerCreate(t *testing.T) {
	t.Run("Test with malformed request body", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Create)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"malformed_request_body","details":"unexpected end of JSON input"}`).
			End()
	})

	t.Run("Test with invalid request body", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Create)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"invalid_request_body","details":"{\"id\":\"cannot be blank\",\"name\":\"cannot be blank\"}"}`).
			End()
	})

	t.Run("Test with object already exists error", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Create)
		topicService.On("Create", mock.Anything).Return(hammer.ErrTopicAlreadyExists)

		apitest.New().
			Handler(r).
			Post("/topics").
			Body(`{"id": "topic", "name": "Topic"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"topic_already_exists","details":""}`).
			End()
	})

	t.Run("Test with service error", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Create)
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

	t.Run("Test success", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Post("/topics", topicHandler.Create)
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

func TestTopicHandlerList(t *testing.T) {
	topics := []hammer.Topic{hammer.MakeTestTopic(), hammer.MakeTestTopic()}

	t.Run("Test with service error", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Get("/topics", topicHandler.List)
		topicService.On("FindAll", mock.Anything, mock.Anything).Return(topics, errors.New("service_error"))

		apitest.New().
			Handler(r).
			Get("/topics").
			Expect(t).
			Status(http.StatusInternalServerError).
			Body(`{"message":"service_error","details":"service_error"}`).
			End()
	})

	t.Run("Test success", func(t *testing.T) {
		topicService := mocks.TopicService{}
		topicHandler := NewTopicHandler(&topicService)
		r := chi.NewRouter()
		r.Get("/topics", topicHandler.List)
		topicService.On("FindAll", mock.Anything, mock.Anything).Return(topics, nil)

		apitest.New().
			Handler(r).
			Get("/topics").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}
