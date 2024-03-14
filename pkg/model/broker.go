package model

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"wowsan/constants"
	grpcClient "wowsan/pkg/broker/transport"
	pb "wowsan/pkg/proto/broker"
	"wowsan/pkg/simulator"
	grpcSubscriberClient "wowsan/pkg/subscriber/transport"

	"github.com/sirupsen/logrus"
)

type Broker struct {
	logger              *logrus.Logger
	RpcClient           grpcClient.BrokerClient
	RpcSubscriberClient grpcSubscriberClient.SubscriberClient
	Id                  string
	Ip                  string
	Port                string
	Brokers             map[string]*Broker
	Publishers          map[string]*Publisher
	Subscribers         map[string]*Subscriber
	SRT                 []*SubscriptionRoutingTableItem
	PRT                 []*PublicationRoutingTableItem

	// Message queue
	MessageQueue chan *MessageRequest

	// 성능 지표 (평균 큐 대기 시간, 평균 서비스 시간, 평균 메시지 도착 간격)
	QueueTime        time.Duration
	ServiceTime      time.Duration
	LastArrivalTime  time.Time
	InterArrivalTime time.Duration
}

// public func
func NewBroker(id, ip, port string, logger *logrus.Logger) *Broker {
	rpcClient := grpcClient.NewBrokerClient()
	rpcSubscriberClient := grpcSubscriberClient.NewSubscriberClient()
	return &Broker{
		logger:              logger,
		RpcClient:           rpcClient,
		RpcSubscriberClient: rpcSubscriberClient,
		Id:                  id,
		Ip:                  ip,
		Port:                port,
		Publishers:          make(map[string]*Publisher),
		Subscribers:         make(map[string]*Subscriber),
		Brokers:             make(map[string]*Broker),
		// SRT:         make([]*SubscriptionRoutingTableItem, 0),
		MessageQueue:     make(chan *MessageRequest, 1000),
		QueueTime:        0,
		ServiceTime:      0,
		InterArrivalTime: 0,
	}
}

func (b *Broker) AddBroker(id string, ip string, port string) *Broker {
	broker := &Broker{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	b.Brokers[id] = broker
	return broker
}

func (b *Broker) AddPublisher(id string, ip string, port string) *Publisher {
	publisher := &Publisher{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	b.Publishers[id] = publisher
	return publisher
}

func (b *Broker) AddSubscriber(id string, ip string, port string) *Subscriber {
	subscriber := &Subscriber{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	b.Subscribers[id] = subscriber
	return subscriber
}

func (b *Broker) DoMessageQueue() {
	var totalQueueTime time.Duration
	var totalServiceTime time.Duration
	var messageCount int64

	for {
		message := <-b.MessageQueue

		// 큐 대기 시간 측정 종료
		queueTime := time.Since(message.EnqueueTime)
		totalQueueTime += queueTime
		messageCount++
		avgQueueTime := totalQueueTime / time.Duration(messageCount)
		b.QueueTime = avgQueueTime
		log.Printf("Cumulative Average Queue Waiting Time: %v\n", avgQueueTime)

		switch message.MessageType {
		case constants.ADVERTISEMENT:
			// 서비스 시간 측정 시작
			message.EnserviceTime = time.Now()

			// time.Sleep(2 * time.Second)
			randomServiceTime := simulator.GetGaussianFigure(2, 0.5)
			time.Sleep(randomServiceTime)

			b.SendAdvertisement(message)

			// 서비스 시간 측정 종료
			serviceTime := time.Since(message.EnserviceTime)

			totalServiceTime += serviceTime
			avgServiceTime := totalServiceTime / time.Duration(messageCount)
			b.ServiceTime = avgServiceTime
			log.Printf("Cumulative Average Service Time: %v\n", avgServiceTime)
		case constants.SUBSCRIPTION:
			// 서비스 시간 측정 시작
			message.EnserviceTime = time.Now()

			time.Sleep(2 * time.Second)
			b.SendSubscription(message)

			// 서비스 시간 측정 종료
			serviceTime := time.Since(message.EnserviceTime)

			totalServiceTime += serviceTime
			avgServiceTime := totalServiceTime / time.Duration(messageCount)
			b.ServiceTime = avgServiceTime
			log.Printf("Cumulative Average Service Time: %v\n", avgServiceTime)
		case constants.PUBLICATION:
			// 서비스 시간 측정 시작
			message.EnserviceTime = time.Now()

			time.Sleep(2 * time.Second)
			b.SendPublication(message)

			// 서비스 시간 측정 종료
			serviceTime := time.Since(message.EnserviceTime)

			totalServiceTime += serviceTime
			avgServiceTime := totalServiceTime / time.Duration(messageCount)
			b.ServiceTime = avgServiceTime
			log.Printf("Cumulative Average Message Service Time: %v\n", avgServiceTime)
		}
	}
}

func (b *Broker) PushMessageToQueue(msgReq *MessageRequest) {
	var totalInterArrivalTime time.Duration
	var messageCount int64

	// 큐 대기 시간 측정 시작
	msgReq.EnqueueTime = time.Now()

	b.MessageQueue <- msgReq
	messageCount++

	if b.LastArrivalTime.IsZero() {
		b.LastArrivalTime = msgReq.EnqueueTime
	} else {
		b.InterArrivalTime = msgReq.EnqueueTime.Sub(b.LastArrivalTime)
		b.LastArrivalTime = msgReq.EnqueueTime

		totalInterArrivalTime += b.InterArrivalTime
		avgInterArrivalTime := totalInterArrivalTime / time.Duration(messageCount)
		b.InterArrivalTime = avgInterArrivalTime
		log.Printf("Cumulative Average Inter-arrival Time: %v\n", avgInterArrivalTime)
	}

}

// go routine 1
// func (b *Broker) DoAdvertisementQueue() {
// 	for {
// 		select {
// 		case reqSrtItem := <-b.AdvertisementQueue:
// 			b.SendAdvertisement(
// 				reqSrtItem,
// 			)
// 		}
// 	}
// }

// func (b *Broker) PushAdvertisementToQueue(req *AdvertisementRequest) {
// 	b.AdvertisementQueue <- req
// }

func (b *Broker) SendAdvertisement(advReq *MessageRequest) error {
	// srt := b.SRT

	reqSrtItem := NewSRTItem(
		advReq.Id,
		advReq.Ip,
		advReq.Port,
		advReq.Subject,
		advReq.Operator,
		advReq.Value,
		advReq.NodeType,
		advReq.HopCount+1,
		advReq.MessageId,
		advReq.SenderId,
	)

	isExist := false
	isShorter := false

	// 새로운 advertisement가 아닌 경우 (같은 advertisement가 이미 존재하는 경우):
	// 건너온 hop이 더 짧으면 last hop을 추가하고
	// 건너온 hop이 같으면 기존 last hop을 대체함.
	for index, item := range b.SRT {
		fmt.Printf("srt: %s %s %s %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.HopCount)
		fmt.Printf("reqSrtItem: %s %s %s %d\n", reqSrtItem.Advertisement.Subject, reqSrtItem.Advertisement.Operator, reqSrtItem.Advertisement.Value, reqSrtItem.HopCount)
		if MatchingEngineSRT(item, reqSrtItem) {
			if item.HopCount >= reqSrtItem.HopCount {
				if item.HopCount == reqSrtItem.HopCount {
					b.SRT[index].AddLastHop(advReq.Id, advReq.Ip, advReq.Port, advReq.NodeType)
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
				fmt.Printf("[SRT] %s %s %s | %s | %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.LastHop[0].Id, item.HopCount)
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
			fmt.Printf("[SRT] %s %s %s | %s | %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.LastHop[0].Id, item.HopCount)
		}
	}
	fmt.Println("Len SRT: ", len(b.SRT))

	// 새로운 advertisement이거나 더 짧은 hop으로 온 경우, 이웃 브로커들에게 advertisement 전파
	if isExist == false || isShorter == true {
		newRequest := &pb.SendMessageRequest{
			Id:        b.Id,
			Ip:        b.Ip,
			Port:      b.Port,
			Subject:   advReq.Subject,
			Operator:  advReq.Operator,
			Value:     advReq.Value,
			NodeType:  advReq.NodeType,
			HopCount:  reqSrtItem.HopCount,
			MessageId: advReq.MessageId,
			SenderId:  advReq.SenderId,
		}

		// show broker list
		fmt.Println("==Neighboring Brokers==")
		for _, neighbor := range b.Brokers {
			fmt.Printf("%s\n", neighbor.Id)
		}

		for _, neighbor := range b.Brokers {
			// 온 방향으로는 전송하지 않음
			if neighbor.Id == advReq.Id {
				continue
			}
			fmt.Println("--Sending to Neighbor--")
			fmt.Println("From: ", b.Id)
			fmt.Println("To:   ", neighbor.Id)
			fmt.Println("-----------------------")

			// 새로운 요청을 이웃에게 전송
			_, err := b.RpcClient.RPCSendAdvertisement(
				neighbor.Ip,   //remote broker ip
				neighbor.Port, //remote broker port
				newRequest.Id,
				newRequest.Ip,
				newRequest.Port,
				newRequest.Subject,
				newRequest.Operator,
				newRequest.Value,
				constants.BROKER,
				newRequest.HopCount,
				newRequest.MessageId,
				newRequest.SenderId,
			)
			if err != nil {
				log.Fatalf("error: %v", err)
				continue
			}
		}
	}

	return nil
}

// go routine 2
// func (b *Broker) DoSubscriptionQueue() {
// 	for {
// 		select {
// 		case reqPrtItem := <-b.SubscriptionQueue:
// 			b.SendSubscription(
// 				reqPrtItem,
// 			)
// 		}
// 	}
// }

// func (b *Broker) PushSubscriptionToQueue(req *SubscriptionRequest) {
// 	b.SubscriptionQueue <- req
// }

func (b *Broker) SendSubscription(subReq *MessageRequest) error {
	reqPrtItem := NewPRTItem(
		subReq.Id,
		subReq.Ip,
		subReq.Port,
		subReq.Subject,
		subReq.Operator,
		subReq.Value,
		subReq.NodeType,
		subReq.MessageId,
		subReq.SenderId,
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
					Id:        b.Id,
					Ip:        b.Ip,
					Port:      b.Port,
					Subject:   subReq.Subject,
					Operator:  subReq.Operator,
					Value:     subReq.Value,
					NodeType:  subReq.NodeType,
					MessageId: subReq.MessageId,
					SenderId:  subReq.SenderId,
				}

				// Show neighboring brokers list
				fmt.Println("==Neighboring Brokers==")
				for _, lastHop := range item.LastHop {
					fmt.Printf("%s\n", lastHop.Id)
				}

				for _, lastHop := range item.LastHop {
					fmt.Println("--Routing Sub via SRT--")
					fmt.Println("From: ", b.Id)
					fmt.Println("To:   ", lastHop.Id)
					fmt.Println("-----------------------")

					// 해당하는 advertisement를 보낸 publisher에게 도달한 경우: 전달 완료
					// if lastHop.NodeType == constants.PUBLISHER {
					if lastHop.Id == item.Identifier.SenderId {
						fmt.Printf("Subscription reached publisher %s\n", lastHop.Id)
						break
					}
					// 	for _, publisher := range b.Publishers {
					// 		if publisher.ID == lastHop.ID {
					// 			// RPCNotifyPublisher
					// 		}
					// 	}

					// 새로운 요청을 SRT의 last hop 브로커에게 전송
					_, err := b.RpcClient.RPCSendSubscription(
						lastHop.Ip,   //remote broker ip
						lastHop.Port, //remote broker port
						newRequest.Id,
						newRequest.Ip,
						newRequest.Port,
						newRequest.Subject,
						newRequest.Operator,
						newRequest.Value,
						constants.BROKER,
						newRequest.MessageId,
						newRequest.SenderId,
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

// go routine 3
// func (b *Broker) DoPublicationQueue() {
// 	for {
// 		select {
// 		case reqItem := <-b.PublicationQueue:
// 			b.SendPublication(
// 				reqItem,
// 			)
// 		}
// 	}
// }

// func (b *Broker) PushPublicationToQueue(req *PublicationRequest) {
// 	b.PublicationQueue <- req
// }

func (b *Broker) SendPublication(pubReq *MessageRequest) error {
	for _, item := range b.PRT {
		// subscription의 subject와 publication의 subject가 같은 경우:
		// 해당 subscription의 last hop으로 publication 전달함.

		if pubReq.MessageId == "" {
			pubReq.MessageId = item.Identifier.MessageId
		}

		if item.Subscription.Subject == pubReq.Subject {
			// 해당하는 subscription을 보낸 subscriber에게 도달할 때까지 hop-by-hop으로 전달
			// (PRT의 last hop을 따라가면서 전달)

			// Show last hops
			fmt.Println("===Matching Last Hop===")
			for _, lastHop := range item.LastHop {
				fmt.Printf("%s\n", lastHop.Id)
			}

			for _, lastHop := range item.LastHop {
				fmt.Println("--Routing Pub via PRT--")
				fmt.Println("From: ", b.Id)
				fmt.Println("To:   ", lastHop.Id)
				fmt.Println("-----------------------")

				// 해당하는 subscription을 보낸 subscriber에게 도달한 경우: 전달 완료
				// if lastHop.NodeType == constants.SUBSCRIBER {
				if lastHop.Id == item.Identifier.SenderId {
					fmt.Printf("Publication reached subscriber %s\n", lastHop.Id)

					// Notify subscriber
					b.RpcSubscriberClient.RPCReceivePublication(
						lastHop.Ip,
						lastHop.Port,
						pubReq.Subject,
						pubReq.Operator,
						pubReq.Value,
					)
					break
				}

				if pubReq.MessageId == "" {
					pubReq.MessageId = item.Identifier.MessageId
				}

				if pubReq.MessageId == item.Identifier.MessageId {
					// 새로운 요청을 SRT의 last hop 브로커에게 전송
					_, err := b.RpcClient.RPCSendPublication(
						lastHop.Ip,   //remote broker ip
						lastHop.Port, //remote broker port
						b.Id,
						b.Ip,
						b.Port,
						pubReq.Subject,
						pubReq.Operator,
						pubReq.Value,
						constants.BROKER,
						pubReq.MessageId,
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
