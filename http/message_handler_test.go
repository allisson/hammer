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

func TestMessageHandlerCreate(t *testing.T) {
	t.Run("Test with malformed request body", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Post("/messages", messageHandler.Create)

		apitest.New().
			Handler(r).
			Post("/messages").
			Body(`{`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"malformed_request_body","details":"unexpected end of JSON input"}`).
			End()
	})

	t.Run("Test with invalid request body", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Post("/messages", messageHandler.Create)

		apitest.New().
			Handler(r).
			Post("/messages").
			Body(`{}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"invalid_request_body","details":"{\"data\":\"cannot be blank\"}"}`).
			End()
	})

	t.Run("Test with topic does not exists error", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Post("/messages", messageHandler.Create)
		messageService.On("Create", mock.Anything).Return(hammer.ErrTopicDoesNotExists)

		apitest.New().
			Handler(r).
			Post("/messages").
			Body(`{"topic_id": "topic_id", "data": "{}"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			Body(`{"message":"topic_does_not_exists","details":""}`).
			End()
	})

	t.Run("Test with service error", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Post("/messages", messageHandler.Create)
		messageService.On("Create", mock.Anything).Return(errors.New("service_error"))

		apitest.New().
			Handler(r).
			Post("/messages").
			Body(`{"topic_id": "topic_id", "data": "{}"}`).
			Expect(t).
			Status(http.StatusInternalServerError).
			Body(`{"message":"service_error","details":"service_error"}`).
			End()
	})

	t.Run("Test success", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Post("/messages", messageHandler.Create)
		messageService.On("Create", mock.Anything).Return(nil)

		apitest.New().
			Handler(r).
			Post("/messages").
			Body(`{"topic_id": "topic_id", "data": "{}"}`).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})
}

func TestMessageHandlerList(t *testing.T) {
	messages := []hammer.Message{hammer.MakeTestMessage(), hammer.MakeTestMessage()}

	t.Run("Test with service error", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Get("/messages", messageHandler.List)
		messageService.On("FindAll", mock.Anything, mock.Anything).Return(messages, errors.New("service_error"))

		apitest.New().
			Handler(r).
			Get("/messages").
			Expect(t).
			Status(http.StatusInternalServerError).
			Body(`{"message":"service_error","details":"service_error"}`).
			End()
	})

	t.Run("Test success", func(t *testing.T) {
		messageService := mocks.MessageService{}
		messageHandler := NewMessageHandler(&messageService)
		r := chi.NewRouter()
		r.Get("/messages", messageHandler.List)
		messageService.On("FindAll", mock.Anything, mock.Anything).Return(messages, nil)

		apitest.New().
			Handler(r).
			Get("/messages").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}
