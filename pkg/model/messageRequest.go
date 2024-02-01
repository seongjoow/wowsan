package model

type AdvertisementRequest struct {
	Id        string
	Ip        string
	Port      string
	Subject   string
	Operator  string
	Value     string
	NodeType  string
	HopCount  int64
	MessageId string
	SenderId  string
}

type SubscriptionRequest struct {
	Id        string
	Ip        string
	Port      string
	Subject   string
	Operator  string
	Value     string
	NodeType  string
	MessageId string
	SenderId  string
}

type PublicationRequest struct {
	Id       string
	Ip       string
	Port     string
	Subject  string
	Operator string
	Value    string
	NodeType string
}
