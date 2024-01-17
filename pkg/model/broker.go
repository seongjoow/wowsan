package model

import (
	"fmt"
	"log"
	"strconv"
	"wowsan/constants"
	grpcClient "wowsan/pkg/broker/transport"
	pb "wowsan/pkg/proto"
)

type Broker struct {
	RpcClient   grpcClient.BrokerClient
	ID          string
	IP          string
	Port        string
	Brokers     map[string]*Broker
	Publishers  map[string]*Publisher
	Subscribers map[string]*Subscriber
	SRT         []*SubscriptionRoutingTableItem
	PRT         []*PublicationRoutingTableItem
}

// public func
func NewBroker(id, ip, port string) *Broker {
	rpcClient := grpcClient.NewBrokerClient()
	return &Broker{
		RpcClient:   rpcClient,
		ID:          id,
		IP:          ip,
		Port:        port,
		Publishers:  make(map[string]*Publisher),
		Subscribers: make(map[string]*Subscriber),
		Brokers:     make(map[string]*Broker),
		// SRT:         make([]*SubscriptionRoutingTableItem, 0), // 필요함?
	}
}

func (b *Broker) AddBroker(id string, ip string, port string) *Broker {
	broker := &Broker{
		ID:   id,
		IP:   ip,
		Port: port,
	}
	b.Brokers[id] = broker
	return broker
}

func (b *Broker) SendAdvertisement(id string, ip string, port string, subject string, operator string, value string, hopCount int64, nodeType string) error {
	// srt := b.SRT

	reqSrtItem := NewSRTItem(
		subject,
		operator,
		value,
		id,
		ip,
		port,
		hopCount+1,
		nodeType,
	)

	isExist := false
	isShorter := false

	// 새로운 advertisement가 아닌 경우 (같은 advertisement가 이미 존재하는 경우):
	// 건너온 hop이 더 짧으면 last hop을 추가하고
	// 건너온 hop이 같으면 기존 last hop을 대체함.
	for index, item := range b.SRT {
		fmt.Printf("srt: %s %s %s %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.HopCount)
		fmt.Printf("reqSrtItem: %s %s %s %d\n", reqSrtItem.Advertisement.Subject, reqSrtItem.Advertisement.Operator, reqSrtItem.Advertisement.Value, reqSrtItem.HopCount)
		if item.Advertisement.Subject == reqSrtItem.Advertisement.Subject &&
			item.Advertisement.Operator == reqSrtItem.Advertisement.Operator &&
			item.Advertisement.Value == reqSrtItem.Advertisement.Value {
			if item.HopCount >= reqSrtItem.HopCount {
				if item.HopCount == reqSrtItem.HopCount {
					b.SRT[index].AddLastHop(id, ip, port, nodeType)
				}
				if item.HopCount > reqSrtItem.HopCount {
					// reqSrtItem.HopCount += 1
					b.SRT[index] = reqSrtItem // 슬라이스의 인덱스를 사용하여 요소 직접 업데이트 (item은 b.SRT의 각 요소에 대한 복사본이라 원본 b.SRT 슬라이스의 요소가 변경되지 않음)
					isShorter = true
				}
			}
			fmt.Println("Same adv already exists.")
			isExist = true

			fmt.Println("============Updated LastHop in SRT============")
			for _, item := range b.SRT {
				fmt.Printf("[SRT] %s %s %s | %s | %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.LastHop[0].ID, item.HopCount)
			}
		}
	}

	// 새로운 advertisement인 경우: SRT에 추가
	if isExist == false {
		// reqSrtItem.HopCount += 1
		// srt = append(srt, reqSrtItem)
		b.SRT = append(b.SRT, reqSrtItem)
		fmt.Println("============Added New Adv to SRT============")
		for _, item := range b.SRT {
			fmt.Printf("[SRT] %s %s %s | %s | %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.LastHop[0].ID, item.HopCount)
		}
	}
	fmt.Println("Len SRT: ", len(b.SRT))

	// 새로운 advertisement이거나 더 짧은 hop으로 온 경우, 이웃 브로커들에게 advertisement 전파
	if isExist == false || isShorter == true {
		newRequest := &pb.SendMessageRequest{
			Subject:  subject,
			Operator: operator,
			Value:    value,
			Id:       b.ID,
			Ip:       b.IP,
			Port:     b.Port,
			HopCount: reqSrtItem.HopCount,
			NodeType: nodeType,
		}

		// show broker list
		fmt.Println("==Neighboring Brokers==")
		for _, neighbor := range b.Brokers {
			fmt.Printf("%s\n", neighbor.ID)
		}

		for _, neighbor := range b.Brokers {
			// 온 방향으로는 전송하지 않음
			if neighbor.ID == id {
				continue
			}
			fmt.Println("--Sending to Neighbor--")
			fmt.Println("From: ", b.ID)
			fmt.Println("To:   ", neighbor.ID)
			fmt.Println("-----------------------")

			// 새로운 요청을 이웃에게 전송
			_, err := b.RpcClient.RPCSendAdvertisement(
				neighbor.IP,   //remote broker ip
				neighbor.Port, //remote broker port
				newRequest.Subject,
				newRequest.Operator,
				newRequest.Value,
				newRequest.Id,
				newRequest.Ip,
				newRequest.Port,
				newRequest.HopCount,
				constants.BROKER,
			)
			if err != nil {
				log.Fatalf("error: %v", err)
				continue
			}
		}
	}

	return nil
}

func (b *Broker) AddPublisher(id string, ip string, port string) *Publisher {
	publisher := &Publisher{
		ID:   id,
		IP:   ip,
		Port: port,
	}
	b.Publishers[id] = publisher
	return publisher
}

func (b *Broker) SendSubscription(id string, ip string, port string, subject string, operator string, value string, nodeType string) error {
	reqPrtItem := NewPRTItem(
		subject,
		operator,
		value,
		id,
		ip,
		port,
		nodeType,
		// hopCount,
	)

	for _, item := range b.SRT {
		// advertisement에 subscription이 포함되는 경우:
		// PRT에 추가하고, 해당 advertisement의 last hop으로 subscription 전달함.
		if item.Advertisement.Subject == reqPrtItem.Subscription.Subject &&
			item.Advertisement.Operator == reqPrtItem.Subscription.Operator {
			advValue, _ := strconv.ParseFloat(item.Advertisement.Value, 64)
			subValue, _ := strconv.ParseFloat(reqPrtItem.Subscription.Value, 64)

			if (item.Advertisement.Operator == ">" && advValue > subValue) ||
				(item.Advertisement.Operator == ">=" && advValue >= subValue) ||
				(item.Advertisement.Operator == "<" && advValue < subValue) ||
				(item.Advertisement.Operator == "<=" && advValue <= subValue) {
				// PRT에 추가
				b.PRT = append(b.PRT, reqPrtItem)

				// 해당하는 advertisement를 보낸 publisher에게 도달할 때까지 hop-by-hop으로 전달
				// (SRT의 last hop을 따라가면서 전달)
				newRequest := &pb.SendMessageRequest{
					Subject:  subject,
					Operator: operator,
					Value:    value,
					Id:       b.ID,
					Ip:       b.IP,
					Port:     b.Port,
					NodeType: nodeType,
				}

				// show broker list
				fmt.Println("==Neighboring Brokers==")
				for _, lastHop := range item.LastHop {
					fmt.Printf("%s\n", lastHop.ID)
				}

				for _, lastHop := range item.LastHop {
					// 해당하는 advertisement를 보낸 publisher에게 도달한 경우: 전달 완료
					if lastHop.NodeType == constants.PUBLISHER {
						break
					}

					fmt.Println("--Routing Sub via SRT--")
					fmt.Println("From: ", b.ID)
					fmt.Println("To:   ", lastHop.ID)
					fmt.Println("-----------------------")

					// 새로운 요청을 SRT의 last hop 브로커에게 전송
					_, err := b.RpcClient.RPCSendSubscription(
						lastHop.IP,   //remote broker ip
						lastHop.Port, //remote broker port
						newRequest.Subject,
						newRequest.Operator,
						newRequest.Value,
						newRequest.Id,
						newRequest.Ip,
						newRequest.Port,
						constants.BROKER,
					)
					if err != nil {
						log.Fatalf("error: %v", err)
						continue
					}
				}
			}
		}
	}

	return nil
}
