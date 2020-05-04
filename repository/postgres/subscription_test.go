package repository

import (
	"testing"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	t.Run("Test Store new Subscription", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Subscription", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)

		tx, err = th.txFactory.New()
		assert.Nil(t, err)
		subscription.Name = "My Subscription III"
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		subscriptionFromRepo, err := th.subscriptionRepo.Find(subscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, subscription.Name, subscriptionFromRepo.Name)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)

		subscriptionFromRepo, err := th.subscriptionRepo.Find(subscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, subscriptionFromRepo.ID, subscription.ID)
		assert.Equal(t, subscriptionFromRepo.Name, subscription.Name)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription1 := hammer.MakeTestSubscription()
		subscription1.TopicID = topic.ID
		subscription2 := hammer.MakeTestSubscription()
		subscription2.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription1)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription2)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		subscriptions, err := th.subscriptionRepo.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(subscriptions))
	})

	t.Run("Test FindByTopic", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		subscription1 := hammer.MakeTestSubscription()
		subscription1.TopicID = topic.ID
		subscription2 := hammer.MakeTestSubscription()
		subscription2.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription1)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(tx, &subscription2)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		subscriptions, err := th.subscriptionRepo.FindByTopic(topic.ID)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(subscriptions))
	})
}
