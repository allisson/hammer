package hammer

import "context"

// TopicRepository interface
type TopicRepository interface {
	Find(ctx context.Context, id string) (*Topic, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Topic, error)
	Store(ctx context.Context, topic *Topic) error
	Delete(ctx context.Context, id string) error
}

// SubscriptionRepository interface
type SubscriptionRepository interface {
	Find(ctx context.Context, id string) (*Subscription, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Subscription, error)
	Store(ctx context.Context, subscription *Subscription) error
	Delete(ctx context.Context, id string) error
}

// MessageRepository interface
type MessageRepository interface {
	Find(ctx context.Context, id string) (*Message, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Message, error)
	Store(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id string) error
}

// DeliveryRepository interface
type DeliveryRepository interface {
	Find(ctx context.Context, id string) (*Delivery, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*Delivery, error)
	Store(ctx context.Context, delivery *Delivery) error
	Dispatch(ctx context.Context) (*DeliveryAttempt, error)
}

// DeliveryAttemptRepository interface
type DeliveryAttemptRepository interface {
	Find(ctx context.Context, id string) (*DeliveryAttempt, error)
	FindAll(ctx context.Context, findOptions FindOptions) ([]*DeliveryAttempt, error)
	Store(ctx context.Context, deliveryAttempt *DeliveryAttempt) error
}

// MigrationRepository interface
type MigrationRepository interface {
	Run(ctx context.Context) error
}
