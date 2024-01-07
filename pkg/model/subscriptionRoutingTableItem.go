package model

type Advertisement struct {
	Subject  string
	Value    string
	Operator string
}

type SubscriptionRoutingTableItem struct {
	Advertisement *Advertisement
	LastHop       []*LastHop
}

func NewSRTItem() *SubscriptionRoutingTableItem {
	return &SubscriptionRoutingTableItem{
		Advertisement: &Advertisement{},
		LastHop:       []*LastHop{},
	}
}

func (srtItem *SubscriptionRoutingTableItem) SetAdvertisement(subject string, value string, operator string) {
	srtItem.Advertisement.Subject = subject
	srtItem.Advertisement.Value = value
	srtItem.Advertisement.Operator = operator
}

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
