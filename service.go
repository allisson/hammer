package hammer

import "net/http"

// TopicService interface
type TopicService interface {
	Find(id string) (Topic, error)
	FindAll(findOptions FindOptions) ([]Topic, error)
	Create(topic *Topic) error
	Update(topic *Topic) error
	Delete(id string) error
}

// SubscriptionService interface
type SubscriptionService interface {
	Find(id string) (Subscription, error)
	FindAll(findOptions FindOptions) ([]Subscription, error)
	Create(subscription *Subscription) error
	Update(subscription *Subscription) error
	Delete(id string) error
}

// MessageService interface
type MessageService interface {
	Find(id string) (Message, error)
	FindAll(findOptions FindOptions) ([]Message, error)
	Create(message *Message) error
	Delete(id string) error
}

// DeliveryService interface
type DeliveryService interface {
	Find(id string) (Delivery, error)
	FindAll(findOptions FindOptions) ([]Delivery, error)
	FindToDispatch(limit, offset int) ([]string, error)
	Dispatch(delivery *Delivery, httpClient *http.Client) error
}

// DeliveryAttemptService interface
type DeliveryAttemptService interface {
	Find(id string) (DeliveryAttempt, error)
	FindAll(findOptions FindOptions) ([]DeliveryAttempt, error)
}

// WorkerService interface
type WorkerService interface {
	Run() error
	Stop() error
}

// MigrationService interface
type MigrationService interface {
	Run() error
}
