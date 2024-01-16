package model

type Publisher struct {
	ID   string
	IP   string
	Port string
	// Brokers map[string]*Broker
	Broker *Broker
}

func NewPublisher(id string, ip string, port string) *Publisher {
	return &Publisher{
		ID:   id,
		IP:   ip,
		Port: port,
		// Brokers: make(map[string]*Broker),
		Broker: nil,
	}
}

func (p *Publisher) SetBroker(id string, ip string, port string) *Broker {
	broker := &Broker{
		ID:   id,
		IP:   ip,
		Port: port,
	}
	p.Broker = broker
	return broker
}
