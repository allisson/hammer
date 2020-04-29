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

		topic := hammer.MakeTestTopic()
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Message", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		message.Data = "My Data III"
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		messageFromRepo, err := th.messageRepo.Find(message.ID)
		assert.Nil(t, err)
		assert.Equal(t, message.Data, messageFromRepo.Data)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		message := hammer.MakeTestMessage()
		message.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
		messageFromRepo, err := th.messageRepo.Find(message.ID)
		assert.Nil(t, err)
		assert.Equal(t, messageFromRepo.ID, message.ID)
		assert.Equal(t, messageFromRepo.Data, message.Data)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		message1 := hammer.MakeTestMessage()
		message1.TopicID = topic.ID
		message2 := hammer.MakeTestMessage()
		message2.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message1)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message2)
		assert.Nil(t, err)
		messages, err := th.messageRepo.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(messages))
	})

	t.Run("Test FindByTopic", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		message1 := hammer.MakeTestMessage()
		message1.TopicID = topic.ID
		message2 := hammer.MakeTestMessage()
		message2.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message1)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message2)
		assert.Nil(t, err)
		messages, err := th.messageRepo.FindByTopic(topic.ID, 50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(messages))
	})
}
