package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func makeTestMessage() hammer.Message {
	id := string(randonInt())
	return hammer.Message{
		ID:                fmt.Sprintf("Message_%s", id),
		Data:              fmt.Sprintf("data_%s", id),
		CreatedDeliveries: false,
		CreatedAt:         time.Now().UTC(),
	}
}

func TestMessage(t *testing.T) {
	t.Run("Test Store new Message", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		message := makeTestMessage()
		message.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.messageRepo.Store(&message)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Message", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		message := makeTestMessage()
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

		topic := makeTestTopic()
		message := makeTestMessage()
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
}
