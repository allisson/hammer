package hammer

import (
	"context"
)

// TopicService interface
type TopicService interface {
	Find(ctx context.Context, id string) (*Topic, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Topic, error)
	Create(ctx context.Context, topic *Topic) error
	Update(ctx context.Context, topic *Topic) error
	Delete(ctx context.Context, id string) error
}

// SubscriptionService interface
type SubscriptionService interface {
	Find(ctx context.Context, id string) (*Subscription, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Subscription, error)
	Create(ctx context.Context, subscription *Subscription) error
	Update(ctx context.Context, subscription *Subscription) error
	Delete(ctx context.Context, id string) error
}

// MessageService interface
type MessageService interface {
	Find(ctx context.Context, id string) (*Message, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Message, error)
	Create(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id string) error
}

// DeliveryService interface
type DeliveryService interface {
	Find(ctx context.Context, id string) (*Delivery, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Delivery, error)
}

// DeliveryAttemptService interface
type DeliveryAttemptService interface {
	Find(ctx context.Context, id string) (*DeliveryAttempt, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*DeliveryAttempt, error)
}

// WorkerService interface
type WorkerService interface {
	Run(ctx context.Context)
	Stop(ctx context.Context)
}

// MigrationService interface
type MigrationService interface {
	Run(ctx context.Context) error
}
