package hammer

import (
	"errors"
	"regexp"
	"time"

	"github.com/allisson/go-env"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// DeliveryStatusPending represents the delivery pending status
	DeliveryStatusPending = "pending"
	// DeliveryStatusFailed represents the delivery failed status
	DeliveryStatusFailed = "failed"
	// DeliveryStatusCompleted represents the delivery completed status
	DeliveryStatusCompleted = "completed"
)

var (
	idRegex = regexp.MustCompile(`^[\w.+-]+$`)
	// ErrTopicAlreadyExists is used when the topic already exists on repository.
	ErrTopicAlreadyExists = errors.New("topic_already_exists")
	// ErrTopicDoesNotExists is used when the topic does not exists on repository.
	ErrTopicDoesNotExists = errors.New("topic_does_not_exists")
	// ErrSubscriptionAlreadyExists is used when the subscription already exists on repository.
	ErrSubscriptionAlreadyExists = errors.New("subscription_already_exists")
	// ErrSubscriptionDoesNotExists is used when the subscription does not exists on repository.
	ErrSubscriptionDoesNotExists = errors.New("subscription_does_not_exists")
	// DefaultPaginationLimit represents a default pagination limit on resource list
	DefaultPaginationLimit = env.GetInt("HAMMER_DEFAULT_PAGINATION_LIMIT", 25)
	// MaxPaginationLimit represents the max value for pagination limit on resource list
	MaxPaginationLimit = env.GetInt("HAMMER_MAX_PAGINATION_LIMIT", 50)
	// DefaultMaxDeliveryAttempts represents a default max delivery attempts for subscription
	DefaultMaxDeliveryAttempts = env.GetInt("HAMMER_DEFAULT_MAX_DELIVERY_ATTEMPTS", 5)
	// DefaultDeliveryAttemptDelay represents a default attempt delay for subscription
	DefaultDeliveryAttemptDelay = env.GetInt("HAMMER_DEFAULT_DELIVERY_ATTEMPT_DELAY", 60)
	// DefaultDeliveryAttemptTimeout represents a default attempt timeout for subscription
	DefaultDeliveryAttemptTimeout = env.GetInt("HAMMER_DEFAULT_DELIVERY_ATTEMPT_TIMEOUT", 5)
)

// Error data
type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// Topic data
type Topic struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate topic
func (t Topic) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.ID, validation.Required, validation.Match(idRegex)),
		validation.Field(&t.Name, validation.Required),
	)
}

// ListTopicsResponse data
type ListTopicsResponse struct {
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Topics []Topic `json:"topics"`
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
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// Validate subscription
func (s Subscription) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required, validation.Match(idRegex)),
		validation.Field(&s.TopicID, validation.Required, validation.Match(idRegex)),
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.URL, validation.Required, is.URL),
		validation.Field(&s.MaxDeliveryAttempts, validation.Min(1)),
		validation.Field(&s.DeliveryAttemptDelay, validation.Min(1)),
		validation.Field(&s.DeliveryAttemptTimeout, validation.Min(1)),
	)
}

// ListSubscriptionsResponse data
type ListSubscriptionsResponse struct {
	Limit         int            `json:"limit"`
	Offset        int            `json:"offset"`
	Subscriptions []Subscription `json:"subscriptions"`
}

// Message data
type Message struct {
	ID                string    `json:"id" db:"id"`
	TopicID           string    `json:"topic_id" db:"topic_id"`
	Data              string    `json:"data" db:"data"`
	CreatedDeliveries bool      `json:"created_deliveries" db:"created_deliveries"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// Validate message
func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Data, validation.Required),
	)
}

// Delivery data
type Delivery struct {
	ID               string    `json:"id" db:"id"`
	TopicID          string    `json:"topic_id" db:"topic_id"`
	SubscriptionID   string    `json:"subscription_id" db:"subscription_id"`
	MessageID        string    `json:"message_id" db:"message_id"`
	ScheduledAt      time.Time `json:"scheduled_at" db:"scheduled_at"`
	DeliveryAttempts int       `json:"delivery_attempts" db:"delivery_attempts"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
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
