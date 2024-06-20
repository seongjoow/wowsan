package model

import (
	"fmt"
	"wowsan/pkg/broker/utils"
)

type PublicationRoutingTable struct {
	Items []*PublicationRoutingTableItem
}

func NewPRT() *PublicationRoutingTable {
	return &PublicationRoutingTable{
		Items: make([]*PublicationRoutingTableItem, 0),
	}
}

func (prt *PublicationRoutingTable) AddItem(
	item *PublicationRoutingTableItem,
) {
	prt.Items = append(prt.Items, item)
}

func (prt *PublicationRoutingTable) RemoveItem(
	subject string,
) {
	for i, item := range prt.Items {
		if item.Subscription.Subject == subject {
			prt.Items = append(prt.Items[:i], prt.Items[i+1:]...)
			return
		}
	}
}

func (prt *PublicationRoutingTable) GetItem(
	subject string,
) (item *PublicationRoutingTableItem, index int) {
	for i, item := range prt.Items {
		if item.Subscription.Subject == subject {
			return item, i
		}
	}
	return nil, -1
}

// print prt table
func (prt *PublicationRoutingTable) PrintPRT() {
	// Define headers and lengths
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "[]LastHop(port)"}
	columnLengths := []int{15, 10, 10, 15, 15, 35}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("PRT")
	for _, item := range prt.Items {
		var lastHop string
		for i, hop := range item.LastHop {
			if i > 0 {
				lastHop += ", "
			}
			lastHop += fmt.Sprintf("%s", hop.Port)
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

func (prt *PublicationRoutingTable) PrintPRTWhereSubject(
	subject string,
) {
	// Define headers and lengths
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "[]LastHop(port)"}
	columnLengths := []int{15, 10, 10, 15, 15, 35}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("PRT")
	for _, item := range prt.Items {
		if item.Subscription.Subject == subject {
			var lastHop string
			for i, hop := range item.LastHop {
				if i > 0 {
					lastHop += ", "
				}
				lastHop += fmt.Sprintf("%s", hop.Port)
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
	}
	table.PrintTable()
}
