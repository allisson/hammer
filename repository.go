package hammer

// TopicRepository interface
type TopicRepository interface {
	Find(id string) (Topic, error)
	FindAll(limit, offset int) ([]Topic, error)
	Store(topic *Topic) error
}

// SubscriptionRepository interface
type SubscriptionRepository interface {
	Find(id string) (Subscription, error)
	FindAll(limit, offset int) ([]Subscription, error)
	FindByTopic(topicID string) ([]Subscription, error)
	Store(subscription *Subscription) error
}

// MessageRepository interface
type MessageRepository interface {
	Find(id string) (Message, error)
	FindAll(limit, offset int) ([]Message, error)
	FindByTopic(topicID string, limit, offset int) ([]Message, error)
	Store(message *Message) error
}

// DeliveryRepository interface
type DeliveryRepository interface {
	Find(id string) (Delivery, error)
	FindAll(limit, offset int) ([]Delivery, error)
	Store(delivery *Delivery) error
}

// DeliveryAttemptRepository interface
type DeliveryAttemptRepository interface {
	Find(id string) (DeliveryAttempt, error)
	FindAll(limit, offset int) ([]DeliveryAttempt, error)
	Store(deliveryAttempt *DeliveryAttempt) error
}

// LockRepository interface
type LockRepository interface {
	Acquire(name string) error
	Close() error
}
