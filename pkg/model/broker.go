package model

import (
	"fmt"
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

// print prt table
func (b *Broker) PrintPRT() {
	// Define headers and lengths
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "[]LastHop(port, nodeType)"}
	columnLengths := []int{15, 10, 10, 15, 15, 35}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("PRT")
	for _, item := range b.PRT {
		var lastHop string
		for i, hop := range item.LastHop {
			if i > 0 {
				lastHop += ", "
			}
			lastHop += fmt.Sprintf("%s(%s)", hop.Port, hop.NodeType)
		}
		row := []string{
			item.Subscription.Subject,
			item.Subscription.Operator,
			item.Subscription.Value,
			item.Identifier.MessageId,
			item.Identifier.SenderId,
			lastHop,
		}
		table.AddRow(row)
	}
	table.PrintTable()
}

// print srt table
func (b *Broker) PrintSRT() {
	// Define headers and lengths
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "HopCount", "[]LastHop(port, nodeType)"}
	columnLengths := []int{15, 10, 10, 15, 15, 10, 35}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("SRT")
	for _, item := range b.SRT {
		var lastHop string
		for i, hop := range item.LastHop {
			if i > 0 {
				lastHop += ", "
			}
			lastHop += fmt.Sprintf("(%s,%s)", hop.Port, hop.NodeType)
		}
		row := []string{
			item.Advertisement.Subject,
			item.Advertisement.Operator,
			item.Advertisement.Value,
			item.Identifier.MessageId,
			item.Identifier.SenderId,
			fmt.Sprintf("%d", item.HopCount),
			lastHop,
		}
		table.AddRow(row)
	}
	table.PrintTable()
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
