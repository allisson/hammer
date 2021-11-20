package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allisson/hammer"
)

func TestDispatchToURL(t *testing.T) {
	t.Run("Invalid delivery url", func(t *testing.T) {
		delivery := hammer.MakeTestDelivery()
		delivery.URL = "http://localhost:9999"

		dr := dispatchToURL(delivery)
		assert.False(t, dr.Success)
		assert.Equal(t, `Post "http://localhost:9999": dial tcp [::1]:9999: connect: connection refused`, dr.Error)
	})

	t.Run("Invalid response status code", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		delivery := hammer.MakeTestDelivery()
		delivery.URL = httpServer.URL

		dr := dispatchToURL(delivery)
		assert.NotEqual(t, "", dr.Response)
		assert.Equal(t, http.StatusInternalServerError, dr.ResponseStatusCode)
		assert.False(t, dr.Success)
		assert.Equal(t, "", dr.Error)
	})

	t.Run("Valid response status code", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		delivery := hammer.MakeTestDelivery()
		delivery.URL = httpServer.URL

		dr := dispatchToURL(delivery)
		assert.NotEqual(t, "", dr.Response)
		assert.Equal(t, http.StatusOK, dr.ResponseStatusCode)
		assert.True(t, dr.Success)
		assert.Equal(t, "", dr.Error)
	})
}

func TestDelivery(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Store new Delivery", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Delivery", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)

		delivery.Status = "completed"
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
		deliveryFromRepo, err := th.deliveryRepo.Find(ctx, delivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, delivery.Status, deliveryFromRepo.Status)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
		deliveryFromRepo, err := th.deliveryRepo.Find(ctx, delivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, deliveryFromRepo.ID, delivery.ID)
		assert.Equal(t, deliveryFromRepo.Status, delivery.Status)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message1 := hammer.MakeTestMessage()
		message1.TopicID = topic.ID
		message2 := hammer.MakeTestMessage()
		message2.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message1)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message2)
		assert.Nil(t, err)
		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		deliveries, err := th.deliveryRepo.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(deliveries))
	})

	t.Run("Test Dispatch", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		subscription.URL = httpServer.URL
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		deliveryAttempt, err := th.deliveryRepo.Dispatch(ctx)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, deliveryAttempt.ResponseStatusCode)
		assert.NotEqual(t, "", deliveryAttempt.Request)
		assert.NotEqual(t, "", deliveryAttempt.Response)
		assert.True(t, deliveryAttempt.Success)
	})
}
