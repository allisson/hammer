package postgres

import (
	"testing"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryAttempt(t *testing.T) {
	t.Run("Test Store new DeliveryAttempt", func(t *testing.T) {
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
		deliveryAttempt := hammer.MakeTestDeliveryAttempt()
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

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		deliveryAttempt := hammer.MakeTestDeliveryAttempt()
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

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		deliveryAttempt := hammer.MakeTestDeliveryAttempt()
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

	t.Run("Test FindAll", func(t *testing.T) {
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
		deliveryAttempt1 := hammer.MakeTestDeliveryAttempt()
		deliveryAttempt1.DeliveryID = delivery.ID
		deliveryAttempt2 := hammer.MakeTestDeliveryAttempt()
		deliveryAttempt2.DeliveryID = delivery.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(&delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(&deliveryAttempt1)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(&deliveryAttempt2)
		assert.Nil(t, err)
		deliveryAttempts, err := th.deliveryAttemptRepo.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(deliveryAttempts))
	})
}
