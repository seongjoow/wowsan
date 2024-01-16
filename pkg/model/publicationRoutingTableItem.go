package model

type Subscription struct {
	Subject  string
	Operator string
	Value    string
}

type PublicationRoutingTableItem struct {
	Subscription *Subscription
	LastHop      []*LastHop
	// HopCount	int64
}

func NewPRTItem(
	subject string,
	operator string,
	value string,
	id string,
	ip string,
	port string,
	nodeType string,
	// hopCount int64,
) *PublicationRoutingTableItem {
	return &PublicationRoutingTableItem{
		Subscription: &Subscription{
			Subject:  subject,
			Operator: operator,
			Value:    value,
		},
		LastHop: []*LastHop{NewLastHop(id, ip, port, nodeType)},
	}
}

func (prtItem *PublicationRoutingTableItem) AddLastHop(id string, ip string, port string, nodeType string) {
	lastHop := NewLastHop(id, ip, port, nodeType)
	prtItem.LastHop = append(prtItem.LastHop, lastHop)
}
