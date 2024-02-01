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
	Identifier *Identifier
}

func NewPRTItem(
	id string,
	ip string,
	port string,
	subject string,
	operator string,
	value string,
	nodeType string,
	// hopCount int64,
	messageId string,
	subscriberId string,
) *PublicationRoutingTableItem {
	return &PublicationRoutingTableItem{
		Subscription: &Subscription{
			Subject:  subject,
			Operator: operator,
			Value:    value,
		},
		LastHop: []*LastHop{NewLastHop(id, ip, port, nodeType)},
		Identifier: &Identifier{
			MessageId: messageId,
			SenderId:  subscriberId,
		},
	}
}

func (prtItem *PublicationRoutingTableItem) AddLastHop(id string, ip string, port string, nodeType string) {
	lastHop := NewLastHop(id, ip, port, nodeType)
	prtItem.LastHop = append(prtItem.LastHop, lastHop)
}
