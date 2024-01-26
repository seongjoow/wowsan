package model

type AdvertisementRequest struct {
	Id       string
	Ip       string
	Port     string
	Subject  string
	Operator string
	Value    string
	HopCount int64
	NodeType string
}

type SubscriptionRequest struct {
	Id       string
	Ip       string
	Port     string
	Subject  string
	Operator string
	Value    string
	NodeType string
}

// type PublicationRequest struct {
// 	Id       string
// 	Ip       string
// 	Port     string
// 	Subject  string
// 	Operator string
// 	Value    string
// 	NodeType string
// }
