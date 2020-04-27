package hammer

// TopicRepository interface
type TopicRepository interface {
	Find(id string) (Topic, error)
	Store(topic *Topic) error
}

// SubscriptionRepository interface
type SubscriptionRepository interface {
	Find(id string) (Subscription, error)
	FindByTopic(topicID string) ([]Subscription, error)
	Store(subscription *Subscription) error
}

// MessageRepository interface
type MessageRepository interface {
	Find(id string) (Message, error)
	Store(message *Message) error
}

// DeliveryRepository interface
type DeliveryRepository interface {
	Find(id string) (Delivery, error)
	Store(delivery *Delivery) error
}

// DeliveryAttemptRepository interface
type DeliveryAttemptRepository interface {
	Find(id string) (DeliveryAttempt, error)
	Store(deliveryAttempt *DeliveryAttempt) error
}

// LockRepository interface
type LockRepository interface {
	Acquire(name string) error
	Close() error
}
