package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MessageRequest struct {
	Id          string
	Ip          string
	Port        string
	Subject     string
	Operator    string
	Value       string
	NodeType    string
	HopCount    int64
	MessageId   string
	SenderId    string
	MessageType string
	HopId       uuid.UUID
	EnqueueTime time.Time
}

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
