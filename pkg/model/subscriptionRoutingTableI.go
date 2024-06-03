package model

import (
	"fmt"
	"wowsan/pkg/broker/utils"
)

type SubscriptionRoutingTable struct {
	Items []*SubscriptionRoutingTableItem
}

func NewSRT() *SubscriptionRoutingTable {
	return &SubscriptionRoutingTable{
		Items: make([]*SubscriptionRoutingTableItem, 0),
	}
}

func (srt *SubscriptionRoutingTable) AddItem(
	item *SubscriptionRoutingTableItem,
) {
	srt.Items = append(srt.Items, item)
}

func (srt *SubscriptionRoutingTable) GetItem(
	subject string,
) *SubscriptionRoutingTableItem {
	for _, item := range srt.Items {
		if item.Advertisement.Subject == subject {
			return item
		}
	}
	return nil
}

// print srt table
func (srt *SubscriptionRoutingTable) PrintSRT() {
	// Define headers and lengths
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "HopCount", "[]LastHop(port)"}
	columnLengths := []int{15, 10, 10, 15, 15, 10, 35}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("SRT")
	for _, item := range srt.Items {
		var lastHop string
		for i, hop := range item.LastHop {
			if i > 0 {
				lastHop += ", "
			}
			lastHop += fmt.Sprintf("%s", hop.Port)
		}
		row := []string{
			item.Advertisement.Subject,
			item.Advertisement.Operator,
			item.Advertisement.Value,
			item.Identifier.MessageId,
			item.Identifier.SenderId,
			fmt.Sprintf("%d", item.HopCount),
			lastHop,
		}
		table.AddRow(row)
	}
	table.PrintTable()
}

func (srt *SubscriptionRoutingTable) PrintSRTWhereSubject(
	subject string,
) {
	// Define headers and lengths
	columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "HopCount", "[]LastHop(port)"}
	columnLengths := []int{15, 10, 10, 15, 15, 10, 35}
	// Print the table
	table := utils.NewTable(columnHeaders, columnLengths)
	table.SetTitle("SRT")
	for _, item := range srt.Items {
		if item.Advertisement.Subject == subject {
			var lastHop string
			for i, hop := range item.LastHop {
				if i > 0 {
					lastHop += ", "
				}
				lastHop += fmt.Sprintf("%s", hop.Port)
			}
			row := []string{
				item.Advertisement.Subject,
				item.Advertisement.Operator,
				item.Advertisement.Value,
				item.Identifier.MessageId,
				item.Identifier.SenderId,
				fmt.Sprintf("%d", item.HopCount),
				lastHop,
			}
			table.AddRow(row)
		}
	}
	table.PrintTable()
}

// 새로운 advertisement가 아닌 경우 (같은 내용의 advertisement가 이미 존재하는 경우):
// 건너온 hop이 더 짧으면 last hop을 추가하고
// 건너온 hop이 같으면 기존 last hop을 대체함.
// if item.HopCount >= reqSrtItem.HopCount {
// 	if item.HopCount == reqSrtItem.HopCount {
// 		uc.broker.SRTmutex.Lock()
// 		uc.broker.SRT.Items[index].AddLastHop(advReq.Id, advReq.Ip, advReq.Port, advReq.NodeType)
// 		uc.broker.SRTmutex.Unlock()
// 		uc.brokerInfoLogger.GetBrokerInfo(
// 			uc.broker,
// 		)
// 	}
// 	if item.HopCount > reqSrtItem.HopCount {
// 		// reqSrtItem.HopCount += 1
// 		uc.broker.SRTmutex.Lock()
// 		uc.broker.SRT.Items[index] = reqSrtItem // 슬라이스의 인덱스를 사용하여 요소 직접 업데이트 (item은 uc.SRT의 각 요소에 대한 복사본이라 원본 uc.SRT 슬라이스의 요소가 변경되지 않음)
// 		uc.broker.SRTmutex.Unlock()
// 		uc.brokerInfoLogger.GetBrokerInfo(uc.broker)

// 		isShorter = true
// 	}
// }

func (srt *SubscriptionRoutingTable) UpdateLastHop(
	reqSrtItem *SubscriptionRoutingTableItem,
) {
	for index, item := range srt.Items {
		if MatchingEngineSRT(item, reqSrtItem) {
			switch {
			case item.HopCount == reqSrtItem.HopCount:
				srt.Items[index].AddLastHop(reqSrtItem.LastHop[0].Id, reqSrtItem.LastHop[0].Ip, reqSrtItem.LastHop[0].Port, reqSrtItem.LastHop[0].NodeType)
			case item.HopCount > reqSrtItem.HopCount:
				reqSrtItem.HopCount += 1
				srt.Items[index] = reqSrtItem
			case item.HopCount < reqSrtItem.HopCount:
				return
			}
		}
	}
}
