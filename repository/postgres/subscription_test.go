package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/allisson/hammer"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func makeTestSubscription() hammer.Subscription {
	id := string(randonInt())
	return hammer.Subscription{
		ID:                     fmt.Sprintf("Subscription_%s", id),
		Name:                   fmt.Sprintf("My Subscription %s", id),
		URL:                    fmt.Sprintf("https://example.com/%s/", id),
		SecretToken:            fmt.Sprintf("token-%s", id),
		MaxDeliveryAttempts:    1,
		DeliveryAttemptDelay:   10,
		DeliveryAttemptTimeout: 5,
		Active:                 true,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}
}

func TestSubscription(t *testing.T) {
	t.Run("Test Store new Subscription", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		subscription := makeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
	})

	t.Run("Test Store against created Subscription", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		subscription := makeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		subscription.Name = "My Subscription III"
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		subscriptionFromRepo, err := th.subscriptionRepo.Find(subscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, subscription.Name, subscriptionFromRepo.Name)
	})

	t.Run("Test Find", func(t *testing.T) {
		th := newTxnTestHelper()
		defer th.db.Close()

		topic := makeTestTopic()
		subscription := makeTestSubscription()
		subscription.TopicID = topic.ID
		err := th.topicRepo.Store(&topic)
		assert.Nil(t, err)
		err = th.subscriptionRepo.Store(&subscription)
		assert.Nil(t, err)
		subscriptionFromRepo, err := th.subscriptionRepo.Find(subscription.ID)
		assert.Nil(t, err)
		assert.Equal(t, subscriptionFromRepo.ID, subscription.ID)
		assert.Equal(t, subscriptionFromRepo.Name, subscription.Name)
	})
}
