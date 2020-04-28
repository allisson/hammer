package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/allisson/go-env"
	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func init() {
	txdb.Register("pgx", "postgres", env.GetString("DATABASE_URL", ""))
}

type txnTestHelper struct {
	db                  *sqlx.DB
	topicRepo           Topic
	subscriptionRepo    Subscription
	messageRepo         Message
	deliveryRepo        Delivery
	deliveryAttemptRepo DeliveryAttempt
}

func newTxnTestHelper() txnTestHelper {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, _ := sqlx.Open("pgx", cName)
	return txnTestHelper{
		db:                  db,
		topicRepo:           NewTopic(db),
		subscriptionRepo:    NewSubscription(db),
		messageRepo:         NewMessage(db),
		deliveryRepo:        NewDelivery(db),
		deliveryAttemptRepo: NewDeliveryAttempt(db),
	}
}

func TestTopic(t *testing.T) {
	t.Run("Test Store new topic", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created topic", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		topic.Name = "My Topic III"
		err = th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		topicFromRepo, err := th.topicRepo.Find(topic.ID)
		assert.Nil(t, err)
		assert.Equal(t, topic.Name, topicFromRepo.Name)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := hammer.MakeTestTopic()
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		topicFromRepo, err := th.topicRepo.Find(topic.ID)
		assert.Nil(t, err)
		assert.Equal(t, topicFromRepo.ID, topic.ID)
		assert.Equal(t, topicFromRepo.Name, topic.Name)
	})

	t.Run("Test FindAll", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic1 := hammer.MakeTestTopic()
		topic2 := hammer.MakeTestTopic()
		err := th.topicRepo.Store(&topic1)
		assert.Nil(t, err)
		err = th.topicRepo.Store(&topic2)
		assert.Nil(t, err)
		topics, err := th.topicRepo.FindAll(50, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(topics))
	})
}
