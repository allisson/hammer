package repository

import (
	"testing"
	"time"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDelivery(t *testing.T) {
	t.Run("Test Store new Delivery", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Delivery", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)

		tx, err = th.txFactory.New()
		assert.Nil(t, err)
		delivery.Status = "completed"
		err = th.deliveryRepo.Store(tx, &delivery)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		deliveryFromRepo, err := th.deliveryRepo.Find(delivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, delivery.Status, deliveryFromRepo.Status)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery := hammer.MakeTestDelivery()
		delivery.TopicID = topic.ID
		delivery.SubscriptionID = subscription.ID
		delivery.MessageID = message.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		deliveryFromRepo, err := th.deliveryRepo.Find(delivery.ID)
		assert.Nil(t, err)
		assert.Equal(t, deliveryFromRepo.ID, delivery.ID)
		assert.Equal(t, deliveryFromRepo.Status, delivery.Status)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery1 := hammer.MakeTestDelivery()
		delivery1.TopicID = topic.ID
		delivery1.SubscriptionID = subscription.ID
		delivery1.MessageID = message.ID
		delivery2 := hammer.MakeTestDelivery()
		delivery2.TopicID = topic.ID
		delivery2.SubscriptionID = subscription.ID
		delivery2.MessageID = message.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery1)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery2)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		deliveries, err := th.deliveryRepo.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(deliveries))
	})

	t.Run("Test FindToDispatch", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		delivery1 := hammer.MakeTestDelivery()
		delivery1.TopicID = topic.ID
		delivery1.SubscriptionID = subscription.ID
		delivery1.MessageID = message.ID
		delivery2 := hammer.MakeTestDelivery()
		delivery2.TopicID = topic.ID
		delivery2.SubscriptionID = subscription.ID
		delivery2.MessageID = message.ID
		delivery2.ScheduledAt = time.Now().Add(time.Duration(1) * time.Hour)
		delivery3 := hammer.MakeTestDelivery()
		delivery3.TopicID = topic.ID
		delivery3.SubscriptionID = subscription.ID
		delivery3.MessageID = message.ID
		delivery3.Status = hammer.DeliveryStatusFailed
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery1)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery2)
		assert.Nil(t, err)
		err = th.deliveryRepo.Store(tx, &delivery3)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		deliveries, err := th.deliveryRepo.FindToDispatch(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deliveries))
		assert.Equal(t, delivery1.ID, deliveries[0].ID)
	})
}
