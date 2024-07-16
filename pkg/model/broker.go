package model

import (
	"sync"
	"time"
	"wowsan/pkg/broker/utils"
)

type Broker struct {
	Id          string
	Ip          string
	Port        string
	Brokers     map[string]*Broker
	Publishers  map[string]*Publisher
	Subscribers map[string]*Subscriber
	// SRT         []*SubscriptionRoutingTableItem
	// PRT         []*PublicationRoutingTableItem
	SRT      *SubscriptionRoutingTable
	PRT      *PublicationRoutingTable
	SRTmutex *sync.Mutex
	PRTmutex *sync.Mutex
	// Message queue
	MessageQueue chan *MessageRequest

	// 성능 지표 (평균 큐 대기 시간, 평균 서비스 시간, 평균 메시지 도착 간격)
	QueueTime               time.Duration
	AverageQueueTime        time.Duration
	ServiceTime             time.Duration
	AverageServiceTime      time.Duration
	LastArrivalTime         time.Time
	InterArrivalTime        time.Duration
	AverageInterArrivalTime time.Duration
	TotalInterArrivalTime   time.Duration
	MessageCount            int

	// 프로그램 종료 여부를 판단하기 위해 메시지 처리 시작 시간을 기록하는 변수
	Close time.Time
}

// public func
func NewBroker(id, ip, port string) *Broker {
	return &Broker{
		Id:                      id,
		Ip:                      ip,
		Port:                    port,
		Publishers:              make(map[string]*Publisher),
		Subscribers:             make(map[string]*Subscriber),
		Brokers:                 make(map[string]*Broker),
		SRT:                     NewSRT(),
		PRT:                     NewPRT(),
		SRTmutex:                new(sync.Mutex),
		PRTmutex:                new(sync.Mutex),
		MessageQueue:            make(chan *MessageRequest, 1000000),
		QueueTime:               0,
		ServiceTime:             0,
		InterArrivalTime:        0,
		AverageInterArrivalTime: 0,
		TotalInterArrivalTime:   0,
		MessageCount:            0,
		Close:                   time.Now(),
	}
}

func (b *Broker) PrintPublisher() {
	// Define headers and lengths
	columnHeaders := []string{"Id", "Ip", "Port"}
	columnLengths := []int{15, 15, 15}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("Publisher")
	for _, publisher := range b.Publishers {
		row := []string{
			publisher.Id,
			publisher.Ip,
			publisher.Port,
		}
		table.AddRow(row)
	}
	table.PrintTable()
}

func (b *Broker) PrintSubscriber() {
	// Define headers and lengths
	columnHeaders := []string{"Id", "Ip", "Port"}
	columnLengths := []int{15, 15, 15}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("Subscriber")
	for _, subscriber := range b.Subscribers {
		row := []string{
			subscriber.Id,
			subscriber.Ip,
			subscriber.Port,
		}
		table.AddRow(row)
	}
	table.PrintTable()
}
