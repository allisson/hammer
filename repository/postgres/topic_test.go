package repository

import (
	"context"
	"testing"

	"github.com/allisson/hammer"
	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	ctx := context.Background()

	t.Run("Test Store new topic", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created topic", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)

		topic.Name = "My Topic III"
		err = th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		topicFromRepo, err := th.topicRepo.Find(ctx, topic.ID)
		assert.Nil(t, err)
		assert.Equal(t, topic.Name, topicFromRepo.Name)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)
		topicFromRepo, err := th.topicRepo.Find(ctx, topic.ID)
		assert.Nil(t, err)
		assert.Equal(t, topicFromRepo.ID, topic.ID)
		assert.Equal(t, topicFromRepo.Name, topic.Name)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic1 := hammer.MakeTestTopic()
		topic2 := hammer.MakeTestTopic()
		err := th.topicRepo.Store(ctx, topic1)
		assert.Nil(t, err)
		err = th.topicRepo.Store(ctx, topic2)
		assert.Nil(t, err)
		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		topics, err := th.topicRepo.FindAll(ctx, findOptions)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(topics))
	})

	t.Run("Test Delete", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(ctx, topic)
		assert.Nil(t, err)

		err = th.topicRepo.Delete(ctx, topic.ID)
		assert.Nil(t, err)
	})
}
