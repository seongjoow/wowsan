package model

import (
	"sync"
	"time"
)

type Broker struct {
	Id          string
	Ip          string
	Port        string
	Brokers     map[string]*Broker
	Publishers  map[string]*Publisher
	Subscribers map[string]*Subscriber
	SRT         []*SubscriptionRoutingTableItem
	PRT         []*PublicationRoutingTableItem
	SRTmutex    *sync.Mutex
	PRTmutex    *sync.Mutex
	// Message queue
	MessageQueue chan *MessageRequest

	// 성능 지표 (평균 큐 대기 시간, 평균 서비스 시간, 평균 메시지 도착 간격)
	QueueTime        time.Duration
	ServiceTime      time.Duration
	LastArrivalTime  time.Time
	InterArrivalTime time.Duration

	// 프로그램 종료 여부를 판단하기 위해 메시지 처리 시작 시간을 기록하는 변수
	Close time.Time
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
		// PRT:         make([]*PublicationRoutingTableItem, 0),
		SRTmutex:         new(sync.Mutex),
		PRTmutex:         new(sync.Mutex),
		MessageQueue:     make(chan *MessageRequest, 1000),
		QueueTime:        0,
		ServiceTime:      0,
		InterArrivalTime: 0,
		Close:            time.Now(),
	}
}
