package hammer

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randonInt() int {
	// nolint:gosec
	return rand.Intn(9999)
}

// MakeTestTopic returns a new Topic
func MakeTestTopic() *Topic {
	id := fmt.Sprintf("%d", randonInt())
	return &Topic{
		ID:        fmt.Sprintf("topic_%s", id),
		Name:      fmt.Sprintf("My Topic %s", id),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// MakeTestSubscription returns a new Subscription
func MakeTestSubscription() *Subscription {
	id := fmt.Sprintf("%d", randonInt())
	return &Subscription{
		ID:                     fmt.Sprintf("Subscription_%s", id),
		Name:                   fmt.Sprintf("My Subscription %s", id),
		URL:                    fmt.Sprintf("https://example.com/%s/", id),
		SecretToken:            fmt.Sprintf("token-%s", id),
		MaxDeliveryAttempts:    5,
		DeliveryAttemptDelay:   60,
		DeliveryAttemptTimeout: 5,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}
}

// MakeTestMessage returns a new Message
func MakeTestMessage() *Message {
	id := fmt.Sprintf("%d", randonInt())
	return &Message{
		ID:          fmt.Sprintf("Message_%s", id),
		ContentType: "application/json",
		Data:        `{"id": "id", "name": "Allisson"}`,
		CreatedAt:   time.Now().UTC(),
	}
}

// MakeTestDelivery returns a new Delivery
func MakeTestDelivery() *Delivery {
	id := fmt.Sprintf("%d", randonInt())
	return &Delivery{
		ID:                     fmt.Sprintf("Delivery_%s", id),
		ContentType:            "application/json",
		Data:                   fmt.Sprintf("data_%s", id),
		URL:                    fmt.Sprintf("https://example.com/%s/", id),
		SecretToken:            fmt.Sprintf("token-%s", id),
		MaxDeliveryAttempts:    5,
		DeliveryAttemptDelay:   60,
		DeliveryAttemptTimeout: 5,
		Status:                 "pending",
		ScheduledAt:            time.Now().UTC(),
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}
}

// MakeTestDeliveryAttempt returns a new DeliveryAttempt
func MakeTestDeliveryAttempt() *DeliveryAttempt {
	id := fmt.Sprintf("%d", randonInt())
	return &DeliveryAttempt{
		ID:                 fmt.Sprintf("DeliveryAttempt_%s", id),
		Success:            true,
		ResponseStatusCode: 201,
		ExecutionDuration:  1000,
		CreatedAt:          time.Now().UTC(),
	}
}

// GenerateULID returns a new ulid id
func GenerateULID() (string, error) {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	// nolint:gosec
	entropy := rand.New(source)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	return id.String(), err
}
