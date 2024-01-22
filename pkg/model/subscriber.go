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

func (s *Subscriber) SetBroker(id string, ip string, port string) *Broker {
	broker := &Broker{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	s.Broker = broker
	return broker
}

func (s *Subscriber) ReceivePublication(id string, subject string, operator string, value string) error {
	//
	return nil
}
