package hammer

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randonInt() int {
	return rand.Intn(9999)
}

// MakeTestTopic returns a new Topic
func MakeTestTopic() Topic {
	id := fmt.Sprintf("%d", randonInt())
	return Topic{
		ID:        fmt.Sprintf("topic_%s", id),
		Name:      fmt.Sprintf("My Topic %s", id),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// MakeTestSubscription returns a new Subscription
func MakeTestSubscription() Subscription {
	id := fmt.Sprintf("%d", randonInt())
	return Subscription{
		ID:                     fmt.Sprintf("Subscription_%s", id),
		Name:                   fmt.Sprintf("My Subscription %s", id),
		URL:                    fmt.Sprintf("https://example.com/%s/", id),
		SecretToken:            fmt.Sprintf("token-%s", id),
		MaxDeliveryAttempts:    1,
		DeliveryAttemptDelay:   10,
		DeliveryAttemptTimeout: 5,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}
}

// MakeTestMessage returns a new Message
func MakeTestMessage() Message {
	id := fmt.Sprintf("%d", randonInt())
	return Message{
		ID:                fmt.Sprintf("Message_%s", id),
		Data:              fmt.Sprintf("data_%s", id),
		CreatedDeliveries: false,
		CreatedAt:         time.Now().UTC(),
	}
}

// MakeTestDelivery returns a new Delivery
func MakeTestDelivery() Delivery {
	id := fmt.Sprintf("%d", randonInt())
	return Delivery{
		ID:          fmt.Sprintf("Delivery_%s", id),
		Status:      "pending",
		ScheduledAt: time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

// MakeTestDeliveryAttempt returns a new DeliveryAttempt
func MakeTestDeliveryAttempt() DeliveryAttempt {
	id := fmt.Sprintf("%d", randonInt())
	return DeliveryAttempt{
		ID:        fmt.Sprintf("DeliveryAttempt_%s", id),
		URL:       fmt.Sprintf("https://example.com/%s/", id),
		Success:   false,
		CreatedAt: time.Now().UTC(),
	}
}
