package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MessageRequest struct {
	Id              string
	Ip              string
	Port            string
	Subject         string
	Operator        string
	Value           string
	NodeType        string
	HopCount        int64
	MessageId       string
	SenderId        string
	MessageType     string
	HopId           uuid.UUID
	EnqueueTime     time.Time
	EnserviceTime   time.Time
	PerformanceInfo []*PerformanceInfo
}

type PerformanceInfo struct {
	BrokerId           string
	Cpu                string
	Memory             string
	QueueLength        string
	QueueTime          string
	AverageQueueTime   string
	ServiceTime        string
	AverageServiceTime string
	ResponseTime       string
	InterArrivalTime   string
	Throughput         string
	AverageThroughput  string
	Timestamp          string
}

// func Test() {
// 	a := &MessageRequest{
// 		Id:         "id",
// 		Ip:         "ip",
// 		Port:       "port",
// 		LogMessage: []LogMessage{},
// 	}
// 	a.LogMessage = append(a.LogMessage, LogMessage{brokerId: "brokerId"})
// 	a.LogMessage = append(a.LogMessage, LogMessage{brokerId: "brokerId2"})
// }

// type AdvertisementRequest struct {
// 	Id        string
// 	Ip        string
// 	Port      string
// 	Subject   string
// 	Operator  string
// 	Value     string
// 	NodeType  string
// 	HopCount  int64
// 	MessageId string
// 	SenderId  string
// 	HopId     uuid.UUID
// }

// type SubscriptionRequest struct {
// 	Id        string
// 	Ip        string
// 	Port      string
// 	Subject   string
// 	Operator  string
// 	Value     string
// 	NodeType  string
// 	MessageId string
// 	SenderId  string
// }

// type PublicationRequest struct {
// 	Id       string
// 	Ip       string
// 	Port     string
// 	Subject  string
// 	Operator string
// 	Value    string
// 	NodeType string
// }
