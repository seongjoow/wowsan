package usecase

import (
	"fmt"
	"log"
	"math/rand"
	// "strconv"
	"time"
	"wowsan/constants"
	brokerTransport "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/broker/utils"
	"wowsan/pkg/logger"
	model "wowsan/pkg/model"

	// "wowsan/pkg/simulator"
	subscriberTransport "wowsan/pkg/subscriber/grpc/client"

	pb "wowsan/pkg/proto/broker"

	"github.com/sirupsen/logrus"
)

type BrokerUsecase interface {
	GetBroker() *model.Broker
	AddBroker(id string, ip string, port string) (*model.Broker, error)
	AddPublisher(id string, ip string, port string) (*model.Publisher, error)
	AddSubscriber(id string, ip string, port string) (*model.Subscriber, error)
	DoMessageQueue()
	PushMessageToQueue(msgReq *model.MessageRequest)
	PerformanceTickLogger(time time.Duration)
}

type brokerUsecase struct {
	hopLogger        *logrus.Logger
	tickLogger       *logrus.Logger
	brokerInfoLogger *logger.BrokerInfoLogger
	broker           *model.Broker
	brokerClient     brokerTransport.BrokerClient
	subscriberClient subscriberTransport.SubscriberClient
	// brokerRPCClient *mo
}

func NewBrokerUsecase(
	hopLogger *logrus.Logger,
	tickLogger *logrus.Logger,
	brokerInfoLogger *logger.BrokerInfoLogger,
	broker *model.Broker,
	brokerClient brokerTransport.BrokerClient,
	subscriberClient subscriberTransport.SubscriberClient,
) BrokerUsecase {
	return &brokerUsecase{
		hopLogger:        hopLogger,
		tickLogger:       tickLogger,
		brokerInfoLogger: brokerInfoLogger,
		broker:           broker,
		brokerClient:     brokerClient,
		subscriberClient: subscriberClient,
	}
}

func (uc *brokerUsecase) GetBroker() *model.Broker {
	return uc.broker
}

func (uc *brokerUsecase) AddBroker(id string, ip string, port string) (*model.Broker, error) {
	broker := &model.Broker{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	uc.broker.Brokers[id] = broker
	return uc.broker, nil
}

func (uc *brokerUsecase) AddPublisher(id string, ip string, port string) (*model.Publisher, error) {
	uc.broker.Publishers[id] = model.NewPublisher(id, ip, port)
	return uc.broker.Publishers[id], nil
}

func (uc *brokerUsecase) AddSubscriber(id string, ip string, port string) (*model.Subscriber, error) {
	uc.broker.Subscribers[id] = model.NewSubscriber(id, ip, port)
	return uc.broker.Subscribers[id], nil
}

func (uc *brokerUsecase) DoMessageQueue() {
	var totalQueueTime time.Duration
	var totalServiceTime time.Duration
	var messageCount int64
	broker := uc.broker
	for {
		message := <-broker.MessageQueue
		message.EnserviceTime = time.Now()
		// time.Sleep(1 * time.Second)
		if broker.Port == "50003" {
			// sleep random time 100 ms ~ 500 ms
			time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
		} else {
			// time.Sleep(10 * time.Second)
			time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
		}
		// 프로그램 종료 여부를 판단하기 위해 메시지 처리 시작 시간을 기록
		uc.broker.Close = time.Now()

		// 큐 대기 시간 측정 종료
		queueTime := time.Since(message.EnqueueTime)
		totalQueueTime += queueTime
		messageCount++
		avgQueueTime := totalQueueTime / time.Duration(messageCount)
		broker.QueueTime = avgQueueTime
		// uc.logger.Printf("Cumulative Average Queue Waiting Time: %v\n", avgQueueTime)
		// log.Printf("Cumulative Average Queue Waiting Time: %v\n", avgQueueTime)

		switch message.MessageType {
		case constants.ADVERTISEMENT:
			// 서비스 시간 측정 시작
			// randomServiceTime := simulator.GetGaussianFigure(2.8, 0.5)
			// time.Sleep(randomServiceTime)
			uc.SendAdvertisement(message)
			// log.Printf("Cumulative Average Service Time: %v\n", avgServiceTime)
		case constants.SUBSCRIPTION:
			// 서비스 시간 측정 시작
			// randomServiceTime := simulator.GetGaussianFigure(0.8, 0.5)
			// time.Sleep(randomServiceTime)
			uc.SendSubscription(message)
			// log.Printf("Cumulative Average Service Time: %v\n", avgServiceTime)

		case constants.PUBLICATION:
			// 서비스 시간 측정 시작
			// randomServiceTime := simulator.GetGaussianFigure(0.8, 0.5)
			// time.Sleep(randomServiceTime)
			uc.SendPublication(message)
			// log.Printf("Cumulative Average Message Service Time: %v\n", avgServiceTime)
		}
		serviceTime := time.Since(message.EnserviceTime)
		totalServiceTime += serviceTime
		avgServiceTime := totalServiceTime / time.Duration(messageCount)
		uc.broker.ServiceTime = avgServiceTime
	}
}

func (uc *brokerUsecase) PushMessageToQueue(msgReq *model.MessageRequest) {
	uc.broker.SRTmutex.Lock()
	defer uc.broker.SRTmutex.Unlock()

	// 큐 대기 시간 측정 시작
	msgReq.EnqueueTime = time.Now()

	uc.broker.MessageQueue <- msgReq
	uc.broker.MessageCount++

	if uc.broker.LastArrivalTime.IsZero() {
		uc.broker.LastArrivalTime = msgReq.EnqueueTime
	} else {
		uc.broker.InterArrivalTime = msgReq.EnqueueTime.Sub(uc.broker.LastArrivalTime)
		uc.broker.LastArrivalTime = msgReq.EnqueueTime

		uc.broker.TotalInterArrivalTime += uc.broker.InterArrivalTime

		avgInterArrivalTime := uc.broker.TotalInterArrivalTime / time.Duration(uc.broker.MessageCount)
		uc.broker.AverageInterArrivalTime = avgInterArrivalTime
		// log.Printf("Cumulative Average Inter-arrival Time: %v\n", avgInterArrivalTime)
	}

}

/*
SRT Handling:
If the advertisement(adv) message is not in the SRT table, the advertisement message is added to the SRT table.
If not, check the hop count of the advertisement message in the SRT table,
1. if the request hop count is shorter, update the hop count and last hop of the advertisement message in the SRT table.
2. if the request hop count is the same, add the last hop of the request to the last hop of the advertisement message in the SRT table.
3. if the request hop count is longer, do nothing.

Broadcast Handling:
If the advertisement message is new or the hop count is shorter,
broadcast the advertisement message to neighboring brokers.
If the advertisement message is not new and the hop count is the same, do nothing.
*/
func (uc *brokerUsecase) SendAdvertisement(advReq *model.MessageRequest) error {
	// simulator.IncreaseCpuUsage()
	// simulator.IncreaseMemoryUsage(100)

	// srt := uc.SRT
	newRequestPerformanceInfo := advReq.PerformanceInfo
	newRequestPerformanceInfo = append(newRequestPerformanceInfo, uc.GetPerformanceInfo())
	uc.hopLogger.WithFields(logrus.Fields{
		"Node":            uc.broker.Id,
		"PerformanceInfo": newRequestPerformanceInfo,
	}).Infof("Advertisement(%s %s %s) %s", advReq.Subject, advReq.Operator, advReq.Value, advReq.MessageId)

	reqSrtItem := model.NewSRTItem(
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

	isShorter := false

	// 새로운 advertisement가 아닌 경우 (같은 내용의 advertisement가 이미 존재하는 경우):
	// 건너온 hop이 같으면 last hop에 주소를 추가하고
	// 건너온 hop이 더 짧으면 기존 last hop의 주소를 대체함.
	item, exists := uc.broker.SRT.GetItem(
		reqSrtItem.Advertisement.Subject,
		reqSrtItem.Advertisement.Operator,
		reqSrtItem.Advertisement.Value,
	)

	// SRT Handling
	switch {
	case !exists:
		fmt.Println("==Added New Adv to SRT==")
		uc.broker.SRTmutex.Lock()
		uc.broker.SRT.AddItem(reqSrtItem)
		uc.broker.SRTmutex.Unlock()
	case exists:
		switch {
		case item.HopCount > reqSrtItem.HopCount:
			// 슬라이스의 인덱스를 사용하여 요소 직접 업데이트 (item은 uc.SRT의 각 요소에 대한 복사본이라 원본 uc.SRT 슬라이스의 요소가 변경되지 않음)
			uc.broker.SRTmutex.Lock()
			uc.broker.SRT.ReplaceItem(reqSrtItem)
			uc.broker.SRTmutex.Unlock()
			isShorter = true
		case item.HopCount == reqSrtItem.HopCount:
			fmt.Println("==Updated LastHop in SRT==")
			uc.broker.SRTmutex.Lock()
			item.AddLastHop(advReq.Id, advReq.Ip, advReq.Port, advReq.NodeType)
			uc.broker.SRTmutex.Unlock()
		case item.HopCount < reqSrtItem.HopCount:
			// do nothing
		}
	}

	// Broadcast Handling
	switch {
	case !exists || isShorter:
		newRequest := &pb.SendMessageRequest{
			Id:        uc.broker.Id,
			Ip:        uc.broker.Ip,
			Port:      uc.broker.Port,
			Subject:   advReq.Subject,
			Operator:  advReq.Operator,
			Value:     advReq.Value,
			NodeType:  advReq.NodeType,
			HopCount:  reqSrtItem.HopCount,
			MessageId: advReq.MessageId,
			SenderId:  advReq.SenderId,
		}

		// Show neighboring brokers list
		fmt.Println("==Neighboring Brokers==")
		for _, neighbor := range uc.broker.Brokers {
			fmt.Printf("%s\n", neighbor.Id)
		}

		for _, neighbor := range uc.broker.Brokers {
			// 온 방향으로는 전송하지 않음
			if neighbor.Id == advReq.Id {
				continue
			}
			// 새로운 요청을 이웃에게 전송
			_, err := uc.brokerClient.RPCSendAdvertisement(
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
				newRequestPerformanceInfo,
			)

			if err != nil {
				uc.hopLogger.Fatalf("error: %v", err)
				continue
			}
		}
	}
	return nil
}

/*
PRT Handling:
At the first edge broker,
If there is advertisement message in SRT that contains the subscription message
and the subscription(sub) message is not in the PRT table,
the subscription message is added to the PRT table.
At the other brokers,
If the subscription message is not in the PRT table,
the subscription message is added to the PRT table.

Broadcast Handling:
If the subscription message matches the advertisement message in the SRT table,
the subscription message is sent to the last hop of the advertisement message.
Then the last hop will do the same thing until it reaches the publisher.
*/
func (uc *brokerUsecase) SendSubscription(subReq *model.MessageRequest) error {
	// simulator.IncreaseCpuUsage()
	// simulator.IncreaseMemoryUsage(100)

	newRequestPerformanceInfo := subReq.PerformanceInfo
	newRequestPerformanceInfo = append(newRequestPerformanceInfo, uc.GetPerformanceInfo())
	if len(newRequestPerformanceInfo) > 7 {
		fmt.Println("PerformanceInfo Infinite loop 발생")
		for _, info := range newRequestPerformanceInfo {
			fmt.Println("Broker ID", info.BrokerId)
			fmt.Println("Broker QueueLength", info.QueueLength)
		}
	}

	uc.hopLogger.WithFields(logrus.Fields{
		"Node":            uc.broker.Id,
		"PerformanceInfo": newRequestPerformanceInfo,
	}).Infof("Subscription(%s %s %s) %s", subReq.Subject, subReq.Operator, subReq.Value, subReq.MessageId)

	reqPrtItem := model.NewPRTItem(
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

	matchedSrtItems, err := uc.broker.SRT.FindMatchesItems(
		reqPrtItem.Subscription.Subject,
		reqPrtItem.Subscription.Operator,
		reqPrtItem.Subscription.Value,
	)
	if !err {
		return nil
	}
	// advertisement의 subject와 operator가 subscription의 subject와 operator와 일치하는 경우:
	for _, srtItem := range matchedSrtItems {
		isSameLastHop := false
		isExistInPRT := false
		// 새로운 subscription이 아닌 경우 (같은 내용의 subscription이 이미 존재하는 경우):
		// PRT의 해당 item에 last hop 슬라이스에 없으면 추가, 있으면 추가하지 않음.
		for index, prtItem := range uc.broker.PRT.Items {
			if prtItem.Identifier.MessageId == reqPrtItem.Identifier.MessageId {
				// PRT의 last hop 슬라이스에 주소가 이미 존재하는 경우:
				// last hop 슬라이스 업데이트하지 않음
				for _, lastHop := range prtItem.LastHop {
					// if lastHop.Address == subReq.Address {
					if lastHop.Id == subReq.Id {
						isSameLastHop = true
						isExistInPRT = true
						break
					}
				}

				// PRT의 last hop 슬라이스에 주소가 존재하지 않는 경우:
				// last hop 슬라이스 업데이트함
				if !isSameLastHop {
					uc.broker.PRTmutex.Lock()
					uc.broker.PRT.Items[index].AddLastHop(subReq.Id, subReq.Ip, subReq.Port, subReq.NodeType)
					uc.broker.PRTmutex.Unlock()
					uc.brokerInfoLogger.GetBrokerInfo(uc.broker)

					isExistInPRT = true
					break
				}
			}
		}

		// 새로운 subscription인 경우:
		// PRT에 추가
		if !isExistInPRT {
			uc.broker.PRTmutex.Lock()
			uc.broker.PRT.AddItem(reqPrtItem)
			uc.broker.PRTmutex.Unlock()
			uc.brokerInfoLogger.GetBrokerInfo(uc.broker)
		}

		// 해당하는 advertisement를 보낸 publisher에게 도달할 때까지 hop-by-hop으로 전달
		// (SRT의 last hop을 따라가면서 전달)

		newRequest := &pb.SendMessageRequest{
			Id:        uc.broker.Id,
			Ip:        uc.broker.Ip,
			Port:      uc.broker.Port,
			Subject:   subReq.Subject,
			Operator:  subReq.Operator,
			Value:     subReq.Value,
			NodeType:  subReq.NodeType,
			MessageId: subReq.MessageId,
			SenderId:  subReq.SenderId,
		}

		// // Show neighboring brokers list
		// fmt.Println("==Neighboring Brokers==")
		// for _, neighbor := range uc.broker.Brokers {
		// 	fmt.Printf("%s\n", neighbor.Id)
		// }

		// Show lasthop list
		fmt.Println("==Lasthop Brokers==")
		for _, lastHop := range srtItem.LastHop {
			fmt.Printf("%s\n", lastHop.Id)
		}

		RoutingSubViaSRTTable := utils.NewFromToTable()
		RoutingSubViaSRTTable.SetTitle("Routing Sub via SRT")
		for index, lastHop := range srtItem.LastHop {
			if index == 0 {
				RoutingSubViaSRTTable.PrintTableTitle()
				RoutingSubViaSRTTable.PrintSeparatorLine()
				RoutingSubViaSRTTable.PrintHeader()
				RoutingSubViaSRTTable.PrintSeparatorLine()
			}
			// 해당하는 advertisement를 보낸 publisher에게 도달한 경우: 전달 완료
			// if lastHop.NodeType == constants.PUBLISHER {
			if lastHop.Id == srtItem.Identifier.SenderId {
				fmt.Printf("Subscription reached publisher %s\n", lastHop.Id)
				RoutingSubViaSRTTable.PrintRowFromTo(srtItem.Identifier.SenderId, lastHop.Id, "Skip", "incoming direction")
				break
			}
			if subReq.Port == lastHop.Port {
				RoutingSubViaSRTTable.PrintRowFromTo(uc.broker.Id, lastHop.Id, "Skip", "same port:"+subReq.Port)
				break
			}
			// 	for _, publisher := range uc.Publishers {
			// 		if publisher.ID == lastHop.ID {
			// 			// RPCNotifyPublisher
			// 		}
			// 	}

			// RoutingSubViaSRTTable.PrintRow([]string{uc.broker.Id, lastHop.Id})
			// 새로운 요청을 SRT의 last hop 브로커에게 전송
			_, err := uc.brokerClient.RPCSendSubscription(
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
				newRequestPerformanceInfo,
			)

			if err != nil {
				log.Printf("error: %v", err)
				RoutingSubViaSRTTable.PrintRowFromTo(uc.broker.Id, lastHop.Id, "Error", err.Error())
				continue
			}
			RoutingSubViaSRTTable.PrintRowFromTo(uc.broker.Id, lastHop.Id, "SendSub", "")
		}
		RoutingSubViaSRTTable.PrintSeparatorLine()
	}
	return nil
}

/*
If the publication message matches the subscription message in the PRT table,
the publication message is sent to the last hop of the subscription message.
Then the last hop will do the same thing until it reaches the subscriber.
*/
func (uc *brokerUsecase) SendPublication(pubReq *model.MessageRequest) error {
	// simulator.IncreaseCpuUsage()
	// simulator.IncreaseMemoryUsage(100)

	newRequestPerformanceInfo := pubReq.PerformanceInfo
	newRequestPerformanceInfo = append(newRequestPerformanceInfo, uc.GetPerformanceInfo())

	uc.hopLogger.WithFields(logrus.Fields{
		"Node":            uc.broker.Id,
		"PerformanceInfo": newRequestPerformanceInfo,
	}).Infof("Publication(%s %s %s) %s", pubReq.Subject, pubReq.Operator, pubReq.Value, pubReq.MessageId)

	for _, item := range uc.broker.PRT.Items {
		// subscription의 subject와 publication의 subject가 같은 경우:
		// 해당 subscription의 last hop으로 publication 전달함.

		if pubReq.MessageId == "" {
			pubReq.MessageId = item.Identifier.MessageId
		}

		if item.Subscription.Subject == pubReq.Subject { // TODO: operator, value도 비교해야 함
			// 해당하는 subscription을 보낸 subscriber에게 도달할 때까지 hop-by-hop으로 전달
			// (PRT의 last hop을 따라가면서 전달)

			// Show last hops
			fmt.Println("===Matching Last Hop===")
			for _, lastHop := range item.LastHop {
				fmt.Printf("%s\n", lastHop.Id)
			}

			RoutingPubViaPRT := utils.NewFromToTable()
			RoutingPubViaPRT.SetTitle("Routing Pub via PRT")
			for index, lastHop := range item.LastHop {
				if index == 0 {
					RoutingPubViaPRT.PrintTableTitle()
					RoutingPubViaPRT.PrintSeparatorLine()
					RoutingPubViaPRT.PrintHeader()
					RoutingPubViaPRT.PrintSeparatorLine()
				}
				// 해당하는 subscription을 보낸 subscriber에게 도달한 경우: 전달 완료
				// if lastHop.NodeType == constants.SUBSCRIBER {
				if lastHop.Id == item.Identifier.SenderId {
					fmt.Printf("Publication reached subscriber %s\n", lastHop.Id)
					RoutingPubViaPRT.PrintRowFromTo(uc.broker.Id, lastHop.Id, "Receive", "")
					// Notify subscriber
					uc.subscriberClient.RPCReceivePublication(
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
					_, err := uc.brokerClient.RPCSendPublication(
						lastHop.Ip,   //remote broker ip
						lastHop.Port, //remote broker port
						uc.broker.Id,
						uc.broker.Ip,
						uc.broker.Port,
						pubReq.Subject,
						pubReq.Operator,
						pubReq.Value,
						constants.BROKER,
						pubReq.MessageId,
						newRequestPerformanceInfo,
					)
					if err != nil {
						log.Printf("error: %v", err)
						RoutingPubViaPRT.PrintRowFromTo(uc.broker.Id, lastHop.Id, "SendPublication", err.Error())
						continue
					}
					RoutingPubViaPRT.PrintRowFromTo(uc.broker.Id, lastHop.Id, "SendPub", "")
				}

			}
		}
	}

	return nil
}

func (uc *brokerUsecase) PerformanceTickLogger(interval time.Duration) {
	broker := uc.broker
	ticker := time.NewTicker(interval)
	bottleneck := bottleneckStatus{}

	defer ticker.Stop()
	for {
		<-ticker.C
		cpu, memory := utils.Utilization()
		queueLength := len(broker.MessageQueue)
		// if queueLength > 1 {
		// 	cpu = cpu + float64(queueLength*10)
		// 	// 또는 랜덤으로
		// }
		queueTime := broker.QueueTime
		serviceTime := broker.ServiceTime
		responseTime := queueTime + serviceTime
		var throughput float64
		if queueTime+serviceTime == 0 {
			throughput = 0
		} else {
			throughput = 1e9 / float64(queueTime+serviceTime) // 초 단위의 값을 얻기 위해서는 나노초 값을 초로 변환 (time.Duration은 기본적으로 나노초 단위의 정수값을 가짐)
		}
		avgerageInterArrivalTime := broker.AverageInterArrivalTime

		uc.tickLogger.WithFields(logrus.Fields{
			"Cpu":               cpu,
			"Memory":            memory,
			"QueueLength":       queueLength,
			"QueueTime":         fmt.Sprintf("%f", queueTime.Seconds()*1000),   // 단위: ms, 소수점 아래 6자리
			"ServiceTime":       fmt.Sprintf("%f", serviceTime.Seconds()*1000), // 단위: ms, 소수점 아래 6자리
			"ResponseTime":      fmt.Sprintf("%f", responseTime.Seconds()*1000),
			"Throughput":        fmt.Sprintf("%.6f", throughput),
			"InterArrivalTime":  fmt.Sprintf("%f", avgerageInterArrivalTime.Seconds()*1000), // 단위: ms, 소수점 아래 6자리
			"Bottleneck":        uc.CheckBottleneck(&bottleneck, memory, queueLength, responseTime),
			"TotalArrival Time": fmt.Sprintf("%f", broker.TotalInterArrivalTime.Seconds()*1000),
			"MessageCount":      broker.MessageCount,
		}).Info("Performance Metrics")

		// LogBrokerInfoToLogfile(broker)
	}
}

type bottleneckStatus struct {
	start           time.Time
	preQueueLength  int
	preResponseTime time.Duration
	state           bool
}

func (uc *brokerUsecase) CheckBottleneck(bottleneckStatus *bottleneckStatus, memory uint64, queueLength int, responseTime time.Duration) bool {
	// TODO: 병목 기준을 어떻게 설정할 것인지 결정해야 함
	// 메모리 사용량이 100000000이면서
	// queue length가 증가하고 response time이 증가하는 상태가 n초 이상 지속되면 병목으로 판단
	var n time.Duration = 3
	if queueLength > bottleneckStatus.preQueueLength && responseTime > bottleneckStatus.preResponseTime {
		if bottleneckStatus.start.IsZero() {
			bottleneckStatus.start = time.Now()
		} else {
			if time.Since(bottleneckStatus.start) > n*time.Second {
				bottleneckStatus.state = true
			}
		}
	} else {
		// set the start time to zero and state to false
		bottleneckStatus.start = time.Time{}
		bottleneckStatus.state = false
	}
	bottleneckStatus.preQueueLength = queueLength
	bottleneckStatus.preResponseTime = responseTime
	return bottleneckStatus.state
}

func (uc *brokerUsecase) GetPerformanceInfo() *model.PerformanceInfo {
	broker := uc.broker
	cpu, memory := utils.Utilization()
	queueLength := len(broker.MessageQueue)
	queueTime := broker.QueueTime
	serviceTime := broker.ServiceTime
	var throughput float64
	if queueTime+serviceTime == 0 {
		throughput = 0
	} else {
		throughput = 1e9 / float64(queueTime+serviceTime) // 초 단위의 값을 얻기 위해서는 나노초 값을 초로 변환 (time.Duration은 기본적으로 나노초 단위의 정수값을 가짐)
	}
	interArrivalTime := broker.InterArrivalTime

	performanceInfo := &model.PerformanceInfo{
		BrokerId:         broker.Id,
		Cpu:              fmt.Sprintf("%f", cpu),
		Memory:           fmt.Sprint(memory),
		QueueLength:      fmt.Sprintf("%d", queueLength),
		QueueTime:        fmt.Sprintf("%f", queueTime.Seconds()*1000),   // 단위: ms, 소수점 아래 6자리
		ServiceTime:      fmt.Sprintf("%f", serviceTime.Seconds()*1000), // 단위: ms, 소수점 아래 6자리
		ResponseTime:     fmt.Sprintf("%f", (queueTime+serviceTime).Seconds()*1000),
		InterArrivalTime: fmt.Sprintf("%f", interArrivalTime.Seconds()*1000), // 단위: ms, 소수점 아래 6자리
		Throughput:       fmt.Sprintf("%.6f", throughput),
		// Timestamp:        fmt.Sprint(time.Now().UnixNano()),
		Timestamp: fmt.Sprint(time.Now()),
	}

	return performanceInfo
}
