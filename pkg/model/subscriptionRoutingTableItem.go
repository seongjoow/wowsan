package model

type Advertisement struct {
	Subject  string
	Operator string
	Value    string
}

type Identifier struct {
	MessageId   string
	PublisherId string
	// TimeStamp
}

type SubscriptionRoutingTableItem struct {
	Advertisement *Advertisement
	LastHop       []*LastHop
	HopCount      int64
	Identifier    *Identifier
}

func NewSRTItem(
	id string,
	ip string,
	port string,
	subject string,
	operator string,
	value string,
	nodeType string,
	hopCount int64,
	messageId string,
	publisherId string,
) *SubscriptionRoutingTableItem {
	return &SubscriptionRoutingTableItem{
		Advertisement: &Advertisement{
			Subject:  subject,
			Operator: operator,
			Value:    value,
		},
		LastHop:  []*LastHop{NewLastHop(id, ip, port, nodeType)},
		HopCount: hopCount,
		Identifier: &Identifier{
			MessageId:   messageId,
			PublisherId: publisherId,
		},
	}
}

func (srtItem *SubscriptionRoutingTableItem) AddLastHop(id string, ip string, port string, nodeType string) {
	lastHop := NewLastHop(id, ip, port, nodeType)
	srtItem.LastHop = append(srtItem.LastHop, lastHop)
}

// func NewSRT(advertisement *Advertisement, lasthop []*LastHop) *SubscriptionRoutingTableItem {
// 	return &SubscriptionRoutingTableItem{
// 		advertisement: *advertisement,
// 		LastHop: lasthop,
// 	}
// }
