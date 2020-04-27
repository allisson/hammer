package hammer

// TopicService interface
type TopicService interface {
	Find(id string) (Topic, error)
	FindAll(limit, offset int) ([]Topic, error)
	Create(topic *Topic) error
	Update(topic *Topic) error
}
