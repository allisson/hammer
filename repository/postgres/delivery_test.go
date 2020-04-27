package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func makeTestDelivery() hammer.Delivery {
	id := fmt.Sprintf("%d", randonInt())
	return hammer.Delivery{
		ID:                     fmt.Sprintf("Delivery_%s", id),
		MaxDeliveryAttempts:    1,
		DeliveryAttemptDelay:   1000,
		DeliveryAttemptTimeout: 1000,
		Status:                 "pending",
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}
}

func TestDelivery(t *testing.T) {
	t.Run("Test Store new Delivery", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		subscription := makeTestSubscription()
		subscription.TopicID = topic.ID
		message := makeTestMessage()
		message.TopicID = topic.ID
		delivery := makeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Delivery", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		subscription := makeTestSubscription()
		subscription.TopicID = topic.ID
		message := makeTestMessage()
		message.TopicID = topic.ID
		delivery := makeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		delivery.Status = "completed"
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		deliveryFromRepo, err := th.deliveryRepo.Find(delivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, delivery.Status, deliveryFromRepo.Status)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		subscription := makeTestSubscription()
		subscription.TopicID = topic.ID
		message := makeTestMessage()
		message.TopicID = topic.ID
		delivery := makeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		deliveryFromRepo, err := th.deliveryRepo.Find(delivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, deliveryFromRepo.ID, delivery.ID)
		assert.Equal(t, deliveryFromRepo.Status, delivery.Status)
	})
}
