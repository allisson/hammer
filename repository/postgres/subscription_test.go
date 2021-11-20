package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allisson/hammer"
)

func TestSubscription(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Store new Subscription", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Subscription", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)

		subscription.Name = "My Subscription III"
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)
		subscriptionFromRepo, err := th.subscriptionRepo.Find(ctx, subscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, subscription.Name, subscriptionFromRepo.Name)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)

		subscriptionFromRepo, err := th.subscriptionRepo.Find(ctx, subscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, subscriptionFromRepo.ID, subscription.ID)
		assert.Equal(t, subscriptionFromRepo.Name, subscription.Name)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription1 := hammer.MakeTestSubscription()
		subscription1.TopicID = topic.ID
		subscription2 := hammer.MakeTestSubscription()
		subscription2.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription1)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription2)
		assert.Nil(t, err)
		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		subscriptions, err := th.subscriptionRepo.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(subscriptions))
	})

	t.Run("Test Delete", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		subscription := hammer.MakeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(ctx, subscription)
		assert.Nil(t, err)

		err = th.subscriptionRepo.Delete(ctx, subscription.ID)
		assert.Nil(t, err)
	})
}
