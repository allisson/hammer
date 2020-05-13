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
	// ErrMessageDoesNotExists is used when the message does not exists on repository.
	ErrMessageDoesNotExists = errors.New("message_does_not_exists")
	// ErrDeliveryDoesNotExists is used when the delivery does not exists on repository.
	ErrDeliveryDoesNotExists = errors.New("delivery_does_not_exists")
	// ErrDeliveryAttemptDoesNotExists is used when the delivery attempt does not exists on repository.
	ErrDeliveryAttemptDoesNotExists = errors.New("delivery_attempt_does_not_exists")
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
	// DefaultSecretTokenLength represents a default length for a random string to secret token if it is not informed
	DefaultSecretTokenLength = env.GetInt("HAMMER_DEFAULT_SECRET_TOKEN_LENGTH", 40)
	// WorkerDatabaseDelay represents a delay for database access by workers
	WorkerDatabaseDelay = env.GetInt("HAMMER_WORKER_DATABASE_DELAY", 5)
	// WorkerDefaultFetchLimit represents the default value for fetch limit
	WorkerDefaultFetchLimit = env.GetInt("HAMMER_WORKER_DEFAULT_FETCH_LIMIT", 100)
)

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

// Message data
type Message struct {
	ID        string    `json:"id" db:"id"`
	TopicID   string    `json:"topic_id" db:"topic_id"`
	Data      string    `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate message
func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Data, validation.Required),
	)
}

// Delivery data
type Delivery struct {
	ID                     string    `json:"id" db:"id"`
	TopicID                string    `json:"topic_id" db:"topic_id"`
	SubscriptionID         string    `json:"subscription_id" db:"subscription_id"`
	MessageID              string    `json:"message_id" db:"message_id"`
	Data                   string    `json:"data" db:"data"`
	URL                    string    `json:"url" db:"url"`
	SecretToken            string    `json:"secret_token" db:"secret_token"`
	MaxDeliveryAttempts    int       `json:"max_delivery_attempts" db:"max_delivery_attempts"`
	DeliveryAttemptDelay   int       `json:"delivery_attempt_delay" db:"delivery_attempt_delay"`
	DeliveryAttemptTimeout int       `json:"delivery_attempt_timeout" db:"delivery_attempt_timeout"`
	ScheduledAt            time.Time `json:"scheduled_at" db:"scheduled_at"`
	DeliveryAttempts       int       `json:"delivery_attempts" db:"delivery_attempts"`
	Status                 string    `json:"status" db:"status"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// DeliveryAttempt data
type DeliveryAttempt struct {
	ID                 string    `json:"id" db:"id"`
	DeliveryID         string    `json:"delivery_id" db:"delivery_id"`
	Request            string    `json:"request" db:"request"`
	Response           string    `json:"response" db:"response"`
	ResponseStatusCode int       `json:"response_status_code" db:"response_status_code"`
	ExecutionDuration  int       `json:"execution_duration" db:"execution_duration"`
	Success            bool      `json:"success" db:"success"`
	Error              string    `json:"error" db:"error"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// WebhookMessage data
type WebhookMessage struct {
	ID             string    `json:"id"`
	TopicID        string    `json:"topic_id"`
	SubscriptionID string    `json:"subscription_id"`
	MessageID      string    `json:"message_id"`
	SecretToken    string    `json:"secret_token"`
	Data           string    `json:"data"`
	CreatedAt      time.Time `json:"created_at"`
}

// FindFilter data
type FindFilter struct {
	FieldName string
	Operator  string
	Value     string
}

// FindPagination data
type FindPagination struct {
	Limit  uint
	Offset uint
}

// FindOrderBy data
type FindOrderBy struct {
	FieldName string
}

// FindOptions data
type FindOptions struct {
	FindFilters    []FindFilter
	FindPagination *FindPagination
	FindOrderBy    *FindOrderBy
}
