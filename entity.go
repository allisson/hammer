package hammer

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// Topic data
type Topic struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Subscription data
type Subscription struct {
	ID                     string    `json:"id" db:"id"`
	TopicID                string    `json:"topic_id" db:"topic_id"`
	Name                   string    `json:"name" db:"name"`
	URL                    string    `json:"url" db:"url"`
	SecretToken            string    `json:"secret_token" db:"secret_token"`
	MaxDeliveryAttempts    int       `json:"max_delivery_attempts" db:"max_delivery_attempts"`
	DeliveryAttemptDelay   int       `json:"delivery_attempt_delay" db:"delivery_attempt_delay"`
	DeliveryAttemptTimeout int       `json:"delivery_attempt_timeout" db:"delivery_attempt_timeout"`
	Active                 bool      `json:"active" db:"active"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// Message data
type Message struct {
	ID                string    `json:"id" db:"id"`
	TopicID           string    `json:"topic_id" db:"topic_id"`
	Data              string    `json:"data" db:"data"`
	CreatedDeliveries bool      `json:"created_deliveries" db:"created_deliveries"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// Delivery data
type Delivery struct {
	ID                     string    `json:"id" db:"id"`
	TopicID                string    `json:"topic_id" db:"topic_id"`
	SubscriptionID         string    `json:"subscription_id" db:"subscription_id"`
	MessageID              string    `json:"message_id" db:"message_id"`
	MaxDeliveryAttempts    int       `json:"max_delivery_attempts" db:"max_delivery_attempts"`
	DeliveryAttemptDelay   int       `json:"delivery_attempt_delay" db:"delivery_attempt_delay"`
	DeliveryAttemptTimeout int       `json:"delivery_attempt_timeout" db:"delivery_attempt_timeout"`
	DeliveryAttempts       int       `json:"delivery_attempts" db:"delivery_attempts"`
	LastDeliveryAttempt    null.Time `json:"last_delivery_attempt" db:"last_delivery_attempt"`
	Status                 string    `json:"status" db:"status"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// DeliveryAttempt data
type DeliveryAttempt struct {
	ID                 string    `json:"id" db:"id"`
	DeliveryID         string    `json:"delivery_id" db:"delivery_id"`
	URL                string    `json:"url" db:"url"`
	RequestHeaders     string    `json:"request_headers" db:"request_headers"`
	RequestBody        string    `json:"request_body" db:"request_body"`
	ResponseHeaders    string    `json:"response_headers" db:"response_headers"`
	ResponseBody       string    `json:"response_body" db:"response_body"`
	ResponseStatusCode int       `json:"response_status_code" db:"response_status_code"`
	ExecutionDuration  int       `json:"execution_duration" db:"execution_duration"`
	Success            bool      `json:"success" db:"success"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
