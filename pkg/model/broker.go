package model

import (
	"fmt"
	"log"
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
	// SRTList    []*SRT
}

// type SRT struct {
// 	adv string
// 	LastHop  LastHop
// }

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

func (b *Broker) SendAdvertisement(id string, ip string, port string, subject string, operator string, value string, hopCount int64) error {
	srt := b.SRT

	reqSrtItem := NewSRTItem(
		subject,
		operator,
		value,
		id,
		ip,
		port,
		hopCount,
	)

	isExist := false
	isShorter := false

	// TODO
	// Publisher로부터 몇 hop 건너온 메시지인지 확인하는 로직 필요

	// 같은 advertisement에 대한 last hop이 이미 존재하는 경우,
	// 건너온 hop이 더 짧거나 같은 경우에만 last hop 업데이트
	for _, item := range srt {
		fmt.Printf("srt: %s %s %s %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.HopCount)
		fmt.Printf("reqSrtItem: %s %s %s %d\n", reqSrtItem.Advertisement.Subject, reqSrtItem.Advertisement.Operator, reqSrtItem.Advertisement.Value, reqSrtItem.HopCount)
		if item.Advertisement.Subject == reqSrtItem.Advertisement.Subject &&
			item.Advertisement.Operator == reqSrtItem.Advertisement.Operator &&
			item.Advertisement.Value == reqSrtItem.Advertisement.Value {
			if item.HopCount >= reqSrtItem.HopCount {
				if item.HopCount == reqSrtItem.HopCount {
					item.AddLastHop(id, ip, port)
					break
				}
				if item.HopCount > reqSrtItem.HopCount {
					reqSrtItem.HopCount += 1
					item = reqSrtItem
					isShorter = true
					break
				}
			}
			fmt.Println("already exist")
			isExist = true
		}
	}

	// 동일한 advertisement가 존재하지 않는 경우 (새로운 advertisement인 경우)
	if isExist == false {
		fmt.Println("========================new advertisement========================")
		reqSrtItem.HopCount = reqSrtItem.HopCount + 1
		// srt = append(srt, reqSrtItem)
		b.SRT = append(b.SRT, reqSrtItem)
	}
	fmt.Println("Len SRT2: ", len(srt))
	// b.SRT = srt // 포인터?로 반환하면 필요 없음?

	// 새로운 advertisement가 추가되었거나, 더 짧은 hop으로 온 경우, advertisement 전파
	if isExist == false || isShorter == true {
		// advertisement가 온 곳으로부터 멀어지는 방향으로 이웃 브로커들에게 advertisement 전달

		newRequest := &pb.SendMessageRequest{
			Subject:  subject,
			Operator: operator,
			Value:    value,
			Id:       b.ID,
			Ip:       b.IP,
			Port:     b.Port,
			HopCount: reqSrtItem.HopCount,
		}

		// show broker list
		for _, neighbor := range b.Brokers {
			fmt.Printf("boroker list: %s \n", neighbor.ID)
		}

		for _, neighbor := range b.Brokers {
			// 온 방향으로는 전송하지 않음
			if neighbor.ID == id {
				continue
			}
			fmt.Println("broker: ", b.ID)
			fmt.Println("neighbor: ", neighbor.ID)

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
			)
			if err != nil {
				log.Fatalf("error: %v", err)
				continue
			}
		}

		// 디버깅용
		for _, item := range srt {
			fmt.Printf("SRT: %s %s %s %d\n", item.LastHop[0].ID, item.LastHop[0].IP, item.LastHop[0].Port, item.HopCount)
			fmt.Printf("SRT: %s %s %s %d\n", item.LastHop[0].ID, item.LastHop[0].IP, item.LastHop[0].Port, item.HopCount)
		}
	}

	return nil
}
