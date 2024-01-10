package model

type Advertisement struct {
	Subject  string
	Operator string
	Value    string
}

type SubscriptionRoutingTableItem struct {
	Advertisement *Advertisement
	LastHop       []*LastHop
	HopCount      int64
}

func NewSRTItem(
	subject string,
	operator string,
	value string,
	id string,
	ip string,
	port string,
	hopCount int64,

) *SubscriptionRoutingTableItem {
	return &SubscriptionRoutingTableItem{
		Advertisement: &Advertisement{
			Subject:  subject,
			Operator: operator,
			Value:    value,
		},
		LastHop:  []*LastHop{NewLastHop(id, ip, port)},
		HopCount: hopCount,
	}
}

// func (srtItem *SubscriptionRoutingTableItem) SetAdvertisement(subject string, operator string, value string) {
// 	srtItem.Advertisement.Subject = subject
// 	srtItem.Advertisement.Operator = operator
// 	srtItem.Advertisement.Value = value
// }

func (srtItem *SubscriptionRoutingTableItem) AddLastHop(id string, ip string, port string) {
	lastHop := NewLastHop(id, ip, port)
	srtItem.LastHop = append(srtItem.LastHop, lastHop)
}

// func NewSRT(advertisement *Advertisement, lasthop []*LastHop) *SubscriptionRoutingTableItem {
// 	return &SubscriptionRoutingTableItem{
// 		advertisement: *advertisement,
// 		LastHop: lasthop,
// 	}
// }
