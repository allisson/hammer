package repository

import (
	"testing"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Run("Test Store new Message", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Message", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)

		tx, err = th.txFactory.New()
		assert.Nil(t, err)
		message.Data = "My Data III"
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		messageFromRepo, err := th.messageRepo.Find(message.ID)
		assert.Nil(t, err)
		assert.Equal(t, message.Data, messageFromRepo.Data)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		messageFromRepo, err := th.messageRepo.Find(message.ID)
		assert.Nil(t, err)
		assert.Equal(t, messageFromRepo.ID, message.ID)
		assert.Equal(t, messageFromRepo.Data, message.Data)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		tx, err := th.txFactory.New()
		assert.Nil(t, err)
		topic := hammer.MakeTestTopic()
		message1 := hammer.MakeTestMessage()
		message1.TopicID = topic.ID
		message2 := hammer.MakeTestMessage()
		message2.TopicID = topic.ID
		err = th.topicRepo.Store(tx, &topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message1)
		assert.Nil(t, err)
		err = th.messageRepo.Store(tx, &message2)
		assert.Nil(t, err)
		err = tx.Commit()
		assert.Nil(t, err)
		findOptions := hammer.FindOptions{
			FindPagination: &hammer.FindPagination{
				Limit:  50,
				Offset: 0,
			},
		}
		messages, err := th.messageRepo.FindAll(findOptions)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(messages))
	})
}
