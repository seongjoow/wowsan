package model

type Subscriber struct {
	Id     string
	Ip     string
	Port   string
	Broker *Broker
}

func NewSubscriber(id string, ip string, port string) *Subscriber {
	return &Subscriber{
		Id:     id,
		Ip:     ip,
		Port:   port,
		Broker: nil,
	}
}

func (s *Subscriber) SetBroker(broker *Broker) {
	s.Broker = broker
}
