package model

type Broker struct {
	// RPCBroker   pb.BrokerServiceClient
	ID          string
	IP          string
	Port        string
	Brokers     map[string]*Broker
	Publishers  map[string]*Publisher
	Subscribers map[string]*Subscriber
	SRT         []*SubscriptionRoutingTableItem
	// SRTList    []*SRT
}

// type SRT struct {
// 	adv string
// 	LastHop  LastHop
// }

// type LastHop struct {
// 	id string
// 	ip string
// 	port string
// }

// public func
func NewBroker(id, ip, port string) *Broker {
	return &Broker{
		ID:          id,
		IP:          ip,
		Port:        port,
		Publishers:  make(map[string]*Publisher),
		Subscribers: make(map[string]*Subscriber),
		Brokers:     make(map[string]*Broker),
	}
}

func (b *Broker) AddBroker(id string, ip string, port string) *Broker {
	broker := &Broker{
		ID:   id,
		IP:   ip,
		Port: port,
	}
	b.Brokers[id] = broker
	return broker
}

// func (b *Broker) SendMessageToBroker(message string) {
// 	for _, broker := range b.Brokers {
// 		b.RPCBroker.SendMessage(
// 			&pb.SendMessageRequest{
// 				Message: message,
// 			},
// 		)
// 		// rpc call

// 	}
//  }

// func (b *Broker) ReceiveMessageFromBroker(message string) {
// }
