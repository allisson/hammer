package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allisson/hammer"
)

func TestDeliveryAttempt(t *testing.T) {
	ctx := context.Background()

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
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(ctx, deliveryAttempt)
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
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(ctx, deliveryAttempt)
		assert.Nil(t, err)

		deliveryAttempt.Success = true
		err = th.deliveryAttemptRepo.Store(ctx, deliveryAttempt)
		assert.Nil(t, err)
		deliveryAttemptFromRepo, err := th.deliveryAttemptRepo.Find(ctx, deliveryAttempt.ID)
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
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(ctx, deliveryAttempt)
		assert.Nil(t, err)
		deliveryAttemptFromRepo, err := th.deliveryAttemptRepo.Find(ctx, deliveryAttempt.ID)
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
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(ctx, message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(ctx, delivery)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(ctx, deliveryAttempt1)
		assert.Nil(t, err)
		err = th.deliveryAttemptRepo.Store(ctx, deliveryAttempt2)
		assert.Nil(t, err)
		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		deliveryAttempts, err := th.deliveryAttemptRepo.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(deliveryAttempts))
	})
}
