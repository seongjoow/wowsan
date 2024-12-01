package model

import (
	"fmt"
	"github.com/google/btree"
	"strconv"
	"wowsan/pkg/broker/utils"
)

type valueItem struct {
	item  *SubscriptionRoutingTableItem
	value float64
}

func (v valueItem) Less(than btree.Item) bool {
	return v.value < than.(valueItem).value
}

type SubscriptionRoutingTable struct {
	items map[string]map[string]*btree.BTree // map[subject]map[operator]*btree.BTree
}

func NewSRT() *SubscriptionRoutingTable {
	return &SubscriptionRoutingTable{
		items: make(map[string]map[string]*btree.BTree),
	}
}

func (m *SubscriptionRoutingTable) GetItem(subject, operator, value string) (*SubscriptionRoutingTableItem, bool) {
	fmt.Println("[SRT] Get Item", subject, operator, value)
	if ops, ok := m.items[subject]; ok {
		if tree, ok := ops[operator]; ok {
			valueFloat, _ := strconv.ParseFloat(value, 64)
			if item := tree.Get(valueItem{value: valueFloat}); item != nil {
				return item.(valueItem).item, true
			}
		}
	}
	return nil, false
}

func (m *SubscriptionRoutingTable) ReplaceItem(item *SubscriptionRoutingTableItem) {
	subject := item.Advertisement.Subject
	operator := item.Advertisement.Operator
	value := item.Advertisement.Value
	fmt.Println("[SRT] ReplaceItem", subject, operator, value)

	// Parse the value to float64
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Error parsing value %s to float64: %v\n", value, err)
		return
	}

	// If the subject map doesn't exist, return
	if _, ok := m.items[subject]; !ok {
		fmt.Println("[SRT] Subject map doesn't exist")
		return
	}

	// If the operator map doesn't exist for the subject, return
	if _, ok := m.items[subject][operator]; !ok {
		fmt.Println("[SRT] Operator map doesn't exist for the subject")
		return
	}

	// Add the SubscriptionRoutingTableItem to the B-tree
	vItem := valueItem{item: item, value: valueFloat}

	// If the SubscriptionRoutingTableItem already exists, replace it
	if existingItem := m.items[subject][operator].Get(vItem); existingItem != nil {
		fmt.Println("[SRT] Item already exists, replace it")
	} else {
		// If the SubscriptionRoutingTableItem doesn't exist, return
		fmt.Println("[SRT] Item doesn't exist, return")
	}
}

func (m *SubscriptionRoutingTable) AddItem(item *SubscriptionRoutingTableItem) {
	subject := item.Advertisement.Subject
	operator := item.Advertisement.Operator
	value := item.Advertisement.Value
	fmt.Println("[SRT] Add Item", subject, operator, value)

	// Parse the value to float64
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Error parsing value %s to float64: %v\n", value, err)
		return
	}

	// If the subject map doesn't exist, create a new one
	if _, ok := m.items[subject]; !ok {
		fmt.Println("[SRT] Create new subject map")
		m.items[subject] = make(map[string]*btree.BTree)
	}

	// If the operator map doesn't exist for the subject, create a new one
	if _, ok := m.items[subject][operator]; !ok {
		fmt.Println("[SRT] Create new operator map")
		m.items[subject][operator] = btree.New(2) // Degree of 2 for B-tree
	}

	// Add the SubscriptionRoutingTableItem to the B-tree
	vItem := valueItem{item: item, value: valueFloat}

	// If the SubscriptionRoutingTableItem already exists, error
	if existingItem := m.items[subject][operator].Get(vItem); existingItem != nil {
		fmt.Println("[SRT] Item already exists")
	} else {
		// If the SubscriptionRoutingTableItem doesn't exist, insert it
		fmt.Println("[SRT] Insert new item")
		m.items[subject][operator].ReplaceOrInsert(vItem)
	}

}

func (m *SubscriptionRoutingTable) GetValuesBySubjectOperator(subject, operator string) ([]float64, bool) {
	values := []float64{}
	if ops, ok := m.items[subject]; ok {
		if tree, ok := ops[operator]; ok {
			tree.Ascend(func(i btree.Item) bool {
				values = append(values, i.(valueItem).value)
				return true
			})
			return values, true
		}
	}
	return nil, false
}

func (m *SubscriptionRoutingTable) PrintSRT() {
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "HopCount", "[]LastHop(port)"}
	columnLengths := []int{15, 10, 10, 15, 15, 10, 35}
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("SRT")

	for subject, ops := range m.items {
		for operator, tree := range ops {
			tree.Ascend(func(i btree.Item) bool {
				vItem := i.(valueItem)
				item := vItem.item
				var lastHop string
				for i, hop := range item.LastHop {
					if i > 0 {
						lastHop += ", "
					}
					lastHop += fmt.Sprintf("%s", hop.Port)
				}
				row := []string{
					subject,
					operator,
					fmt.Sprintf("%f", vItem.value),
					item.Identifier.MessageId,
					item.Identifier.SenderId,
					fmt.Sprintf("%d", item.HopCount),
					lastHop,
				}
				table.AddRow(row)
				return true
			})
		}
	}
	table.PrintTable()
}

func (m *SubscriptionRoutingTable) PrintSRTWhereSubject(subject string) {
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "HopCount", "[]LastHop(port)"}
	columnLengths := []int{15, 10, 10, 15, 15, 10, 35}
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("SRT")

	if ops, ok := m.items[subject]; ok {
		for operator, tree := range ops {
			tree.Ascend(func(i btree.Item) bool {
				vItem := i.(valueItem)
				item := vItem.item
				var lastHop string
				for i, hop := range item.LastHop {
					if i > 0 {
						lastHop += ", "
					}
					lastHop += fmt.Sprintf("%s", hop.Port)
				}
				row := []string{
					subject,
					operator,
					fmt.Sprintf("%f", vItem.value),
					item.Identifier.MessageId,
					item.Identifier.SenderId,
					fmt.Sprintf("%d", item.HopCount),
					lastHop,
				}
				table.AddRow(row)
				return true
			})
		}
	}
	table.PrintTable()
}

func (srt *SubscriptionRoutingTable) UpdateLastHop(reqSrtItem *SubscriptionRoutingTableItem) {
	subject := reqSrtItem.Advertisement.Subject
	operator := reqSrtItem.Advertisement.Operator
	value := reqSrtItem.Advertisement.Value
	valueFloat, _ := strconv.ParseFloat(value, 64)
	fmt.Println("UpdateLastHop", subject, operator, value)
	if item, exists := srt.GetItem(subject, operator, value); exists {
		switch {
		case item.HopCount == reqSrtItem.HopCount:
			fmt.Println("UpdateLastHop: HopCount == reqSrtItem.HopCount, inserting item")
			item.AddLastHop(reqSrtItem.LastHop[0].Id, reqSrtItem.LastHop[0].Ip, reqSrtItem.LastHop[0].Port, reqSrtItem.LastHop[0].NodeType)
		case item.HopCount > reqSrtItem.HopCount:
			fmt.Println("UpdateLastHop: HopCount > reqSrtItem.HopCount, replacing item")
			reqSrtItem.HopCount += 1
			srt.items[subject][operator].ReplaceOrInsert(valueItem{item: reqSrtItem, value: valueFloat})
		case item.HopCount < reqSrtItem.HopCount:
			return
		}
	}
}

func (m *SubscriptionRoutingTable) FindMatchesItems(subject, operator, value string) ([]*SubscriptionRoutingTableItem, bool) {
	fmt.Println("FindMatchesItems", subject, operator, value)
	matches := []*SubscriptionRoutingTableItem{}
	subValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Error parsing value %s to float64: %v\n", value, err)
		return nil, false
	}

	if ops, exists := m.items[subject]; exists {
		if tree, exists := ops[operator]; exists {
			switch operator {
			case ">":
				tree.AscendGreaterOrEqual(valueItem{value: subValue}, func(i btree.Item) bool {
					item := i.(valueItem)
					if item.value >= subValue {
						matches = append(matches, item.item)
					}
					return true
				})
			case ">=":
				tree.AscendGreaterOrEqual(valueItem{value: subValue}, func(i btree.Item) bool {
					item := i.(valueItem)
					if item.value >= subValue {
						matches = append(matches, item.item)
					}
					return true
				})
			case "<":
				tree.DescendLessOrEqual(valueItem{value: subValue}, func(i btree.Item) bool {
					item := i.(valueItem)
					if item.value <= subValue {
						matches = append(matches, item.item)
					}
					return true
				})
			case "<=":
				tree.DescendLessOrEqual(valueItem{value: subValue}, func(i btree.Item) bool {
					item := i.(valueItem)
					if item.value <= subValue {
						matches = append(matches, item.item)
					}
					return true
				})
			}

			fmt.Printf("Found %d matches for subject %s, operator %s, value %s\n", len(matches), subject, operator, value)
			return matches, len(matches) > 0
		}
	}
	return nil, false
}
