package model

type Subscriber struct {
	ID    string
	Topic string
}

func NewSubscriber(id string, topic string) *Subscriber {
	return &Subscriber{
		ID:    id,
		Topic: topic,
	}
}
