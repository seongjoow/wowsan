package model

type Publisher struct {
	Id   string
	Ip   string
	Port string
	// Brokers map[string]*Broker
	Broker *Broker
}

func NewPublisher(id string, ip string, port string) *Publisher {
	return &Publisher{
		Id:   id,
		Ip:   ip,
		Port: port,
		// Brokers: make(map[string]*Broker),
		Broker: nil,
	}
}

func (p *Publisher) SetBroker(id string, ip string, port string) *Broker {
	broker := &Broker{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	p.Broker = broker
	return broker
}
