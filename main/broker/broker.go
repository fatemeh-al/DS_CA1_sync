package broker

type Message interface{}

type Subscriber interface {
	CreateChannel(channel string) (<-chan Message, error)
	DeleteChannel(channel string) error
	Subscribe(channel string) (<-chan Message, error)
}

type Publisher interface {
	Publish(channel string, m Message) error
}

type Broker interface {
	Subscriber
	Publisher
	Close() error
}
