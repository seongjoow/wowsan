package model

import "time"

type Broker struct {
	Id          string
	Ip          string
	Port        string
	Brokers     map[string]*Broker
	Publishers  map[string]*Publisher
	Subscribers map[string]*Subscriber
	SRT         []*SubscriptionRoutingTableItem
	PRT         []*PublicationRoutingTableItem
	// Message queue
	MessageQueue chan *MessageRequest

	// 성능 지표 (평균 큐 대기 시간, 평균 서비스 시간, 평균 메시지 도착 간격)
	QueueTime        time.Duration
	ServiceTime      time.Duration
	LastArrivalTime  time.Time
	InterArrivalTime time.Duration
}

// public func
func NewBroker(id, ip, port string) *Broker {
	return &Broker{
		Id:          id,
		Ip:          ip,
		Port:        port,
		Publishers:  make(map[string]*Publisher),
		Subscribers: make(map[string]*Subscriber),
		Brokers:     make(map[string]*Broker),
		// SRT:         make([]*SubscriptionRoutingTableItem, 0),
		MessageQueue:     make(chan *MessageRequest, 1000),
		QueueTime:        0,
		ServiceTime:      0,
		InterArrivalTime: 0,
	}
}
