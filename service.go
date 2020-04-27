package hammer

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
