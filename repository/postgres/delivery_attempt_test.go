package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func makeTestDeliveryAttempt() hammer.DeliveryAttempt {
	id := string(randonInt())
	return hammer.DeliveryAttempt{
		ID:        fmt.Sprintf("DeliveryAttempt_%s", id),
		URL:       fmt.Sprintf("https://example.com/%s/", id),
		Success:   false,
		CreatedAt: time.Now().UTC(),
	}
}

func TestDeliveryAttempt(t *testing.T) {
	t.Run("Test Store new DeliveryAttempt", func(t *testing.T) {
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
		deliveryAttempt := makeTestDeliveryAttempt()
		deliveryAttempt.DeliveryID = delivery.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(&deliveryAttempt)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created DeliveryAttempt", func(t *testing.T) {
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
		deliveryAttempt := makeTestDeliveryAttempt()
		deliveryAttempt.DeliveryID = delivery.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(&deliveryAttempt)
		assert.Nil(t, err)
		deliveryAttempt.Success = true
		err = th.deliveryAttemptRepo.Store(&deliveryAttempt)
		assert.Nil(t, err)
		deliveryAttemptFromRepo, err := th.deliveryAttemptRepo.Find(deliveryAttempt.ID)
		assert.Nil(t, err)
		assert.Equal(t, deliveryAttempt.Success, deliveryAttemptFromRepo.Success)
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
		deliveryAttempt := makeTestDeliveryAttempt()
		deliveryAttempt.DeliveryID = delivery.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(&deliveryAttempt)
		assert.Nil(t, err)
		deliveryAttemptFromRepo, err := th.deliveryAttemptRepo.Find(deliveryAttempt.ID)
		assert.Nil(t, err)
		assert.Equal(t, deliveryAttemptFromRepo.ID, deliveryAttempt.ID)
		assert.Equal(t, deliveryAttemptFromRepo.Success, deliveryAttempt.Success)
	})
}
