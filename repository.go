package hammer

// TopicRepository interface
type TopicRepository interface {
	Find(id string) (Topic, error)
	FindAll(findOptions FindOptions) ([]Topic, error)
	Store(tx TxRepository, topic *Topic) error
	Delete(tx TxRepository, id string) error
}

// SubscriptionRepository interface
type SubscriptionRepository interface {
	Find(id string) (Subscription, error)
	FindAll(findOptions FindOptions) ([]Subscription, error)
	Store(tx TxRepository, subscription *Subscription) error
	Delete(tx TxRepository, id string) error
}

// MessageRepository interface
type MessageRepository interface {
	Find(id string) (Message, error)
	FindAll(findOptions FindOptions) ([]Message, error)
	Store(tx TxRepository, message *Message) error
	Delete(tx TxRepository, id string) error
}

// DeliveryRepository interface
type DeliveryRepository interface {
	Find(id string) (Delivery, error)
	FindAll(findOptions FindOptions) ([]Delivery, error)
	FindToDispatch(limit, offset int) ([]string, error)
	Store(tx TxRepository, delivery *Delivery) error
}

// DeliveryAttemptRepository interface
type DeliveryAttemptRepository interface {
	Find(id string) (DeliveryAttempt, error)
	FindAll(findOptions FindOptions) ([]DeliveryAttempt, error)
	Store(tx TxRepository, deliveryAttempt *DeliveryAttempt) error
}

// TxRepository interface
type TxRepository interface {
	Exec(query string, arg interface{}) error
	Commit() error
	Rollback() error
}

// TxFactoryRepository interface
type TxFactoryRepository interface {
	New() (TxRepository, error)
}

// MigrationRepository interface
type MigrationRepository interface {
	Run() error
}
