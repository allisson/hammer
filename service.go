package hammer

import "net/http"

// TopicService interface
type TopicService interface {
	Find(id string) (Topic, error)
	FindAll(limit, offset int) ([]Topic, error)
	Create(topic *Topic) error
	Update(topic *Topic) error
}

// SubscriptionService interface
type SubscriptionService interface {
	Find(id string) (Subscription, error)
	FindAll(limit, offset int) ([]Subscription, error)
	Create(subscription *Subscription) error
	Update(subscription *Subscription) error
}

// MessageService interface
type MessageService interface {
	Find(id string) (Message, error)
	FindAll(limit, offset int) ([]Message, error)
	FindByTopic(topicID string, limit, offset int) ([]Message, error)
	Create(message *Message) error
}

// DeliveryService interface
type DeliveryService interface {
	Find(id string) (Delivery, error)
	FindAll(limit, offset int) ([]Delivery, error)
	FindToDispatch(limit, offset int) ([]Delivery, error)
	Dispatch(delivery *Delivery, httpClient *http.Client) error
}
