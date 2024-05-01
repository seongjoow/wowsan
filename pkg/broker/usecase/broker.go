package usecase

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"wowsan/constants"
	brokerTransport "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/broker/utils"
	model "wowsan/pkg/model"
	"wowsan/pkg/simulator"
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
	broker           *model.Broker
	brokerClient     brokerTransport.BrokerClient
	subscriberClient subscriberTransport.SubscriberClient
	// brokerRPCClient *mo
}

func NewBrokerUsecase(
	hopLogger *logrus.Logger,
	tickLogger *logrus.Logger,
	broker *model.Broker,
	brokerClient brokerTransport.BrokerClient,
	subscriberClient subscriberTransport.SubscriberClient,
) BrokerUsecase {
	return &brokerUsecase{
		hopLogger:        hopLogger,
		tickLogger:       tickLogger,
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
			message.EnserviceTime = time.Now()

			// time.Sleep(2 * time.Second)
			randomServiceTime := simulator.GetGaussianFigure(0.8, 0.5)
			time.Sleep(randomServiceTime)

			uc.SendAdvertisement(message)

			// 서비스 시간 측정 종료
			serviceTime := time.Since(message.EnserviceTime)

			totalServiceTime += serviceTime
			avgServiceTime := totalServiceTime / time.Duration(messageCount)
			uc.broker.ServiceTime = avgServiceTime

			// log.Printf("Cumulative Average Service Time: %v\n", avgServiceTime)
		case constants.SUBSCRIPTION:
			// 서비스 시간 측정 시작
			message.EnserviceTime = time.Now()

			// time.Sleep(2 * time.Second)
			randomServiceTime := simulator.GetGaussianFigure(0.8, 0.5)
			time.Sleep(randomServiceTime)

			uc.SendSubscription(message)

			// 서비스 시간 측정 종료
			serviceTime := time.Since(message.EnserviceTime)

			totalServiceTime += serviceTime
			avgServiceTime := totalServiceTime / time.Duration(messageCount)
			uc.broker.ServiceTime = avgServiceTime
			// log.Printf("Cumulative Average Service Time: %v\n", avgServiceTime)
		case constants.PUBLICATION:
			// 서비스 시간 측정 시작
			message.EnserviceTime = time.Now()

			// time.Sleep(2 * time.Second)
			randomServiceTime := simulator.GetGaussianFigure(0.8, 0.5)
			time.Sleep(randomServiceTime)

			uc.SendPublication(message)

			// 서비스 시간 측정 종료
			serviceTime := time.Since(message.EnserviceTime)

			totalServiceTime += serviceTime
			avgServiceTime := totalServiceTime / time.Duration(messageCount)
			uc.broker.ServiceTime = avgServiceTime
			// log.Printf("Cumulative Average Message Service Time: %v\n", avgServiceTime)
		}
	}
}

func (uc *brokerUsecase) PushMessageToQueue(msgReq *model.MessageRequest) {
	var totalInterArrivalTime time.Duration
	var messageCount int64

	// 큐 대기 시간 측정 시작
	msgReq.EnqueueTime = time.Now()

	uc.broker.MessageQueue <- msgReq
	messageCount++

	if uc.broker.LastArrivalTime.IsZero() {
		uc.broker.LastArrivalTime = msgReq.EnqueueTime
	} else {
		uc.broker.InterArrivalTime = msgReq.EnqueueTime.Sub(uc.broker.LastArrivalTime)
		uc.broker.LastArrivalTime = msgReq.EnqueueTime

		totalInterArrivalTime += uc.broker.InterArrivalTime
		avgInterArrivalTime := totalInterArrivalTime / time.Duration(messageCount)
		uc.broker.InterArrivalTime = avgInterArrivalTime
		// log.Printf("Cumulative Average Inter-arrival Time: %v\n", avgInterArrivalTime)
	}

}

// go routine 1
// func (uc *brokerUsecase) DoAdvertisementQueue() {
// 	for {
// 		select {
// 		case reqSrtItem := <-uc.AdvertisementQueue:
// 			uc.SendAdvertisement(
// 				reqSrtItem,
// 			)
// 		}
// 	}
// }

// func (uc *brokerUsecase) PushAdvertisementToQueue(req *AdvertisementRequest) {
// 	uc.AdvertisementQueue <- req
// }

func (uc *brokerUsecase) SendAdvertisement(advReq *model.MessageRequest) error {
	// simulator.IncreaseCpuUsage()
	simulator.IncreaseMemoryUsage(100)

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

	isExist := false
	isShorter := false

	// 새로운 advertisement가 아닌 경우 (같은 advertisement가 이미 존재하는 경우):
	// 건너온 hop이 더 짧으면 last hop을 추가하고
	// 건너온 hop이 같으면 기존 last hop을 대체함.
	for index, item := range uc.broker.SRT {
		fmt.Printf("srt: %s %s %s %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.HopCount)
		fmt.Printf("reqSrtItem: %s %s %s %d\n", reqSrtItem.Advertisement.Subject, reqSrtItem.Advertisement.Operator, reqSrtItem.Advertisement.Value, reqSrtItem.HopCount)
		if model.MatchingEngineSRT(item, reqSrtItem) {
			if item.HopCount >= reqSrtItem.HopCount {
				if item.HopCount == reqSrtItem.HopCount {
					uc.broker.SRTmutex.Lock()
					uc.broker.SRT[index].AddLastHop(advReq.Id, advReq.Ip, advReq.Port, advReq.NodeType)
					uc.broker.SRTmutex.Unlock()
				}
				if item.HopCount > reqSrtItem.HopCount {
					// reqSrtItem.HopCount += 1
					uc.broker.SRTmutex.Lock()
					uc.broker.SRT[index] = reqSrtItem // 슬라이스의 인덱스를 사용하여 요소 직접 업데이트 (item은 uc.SRT의 각 요소에 대한 복사본이라 원본 uc.SRT 슬라이스의 요소가 변경되지 않음)
					uc.broker.SRTmutex.Unlock()

					isShorter = true
				}
			}
			fmt.Println("Same adv already exists.")

			isExist = true

			fmt.Println("============Updated LastHop in SRT============")
			for _, item := range uc.broker.SRT {
				fmt.Printf("[SRT] %s %s %s | %s | %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.LastHop[0].Id, item.HopCount)
			}
		}
	}

	// 새로운 advertisement인 경우: SRT에 추가
	if isExist == false {
		// reqSrtItem.HopCount += 1
		// srt = append(srt, reqSrtItem)
		uc.broker.SRTmutex.Lock()
		uc.broker.SRT = append(uc.broker.SRT, reqSrtItem)
		uc.broker.SRTmutex.Unlock()
		fmt.Println("============Added New Adv to SRT============")
		for _, item := range uc.broker.SRT {
			fmt.Printf("[SRT] %s %s %s | %s | %d\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.LastHop[0].Id, item.HopCount)
		}
	}
	fmt.Println("Len SRT: ", len(uc.broker.SRT))

	// 새로운 advertisement이거나 더 짧은 hop으로 온 경우, 이웃 브로커들에게 advertisement 전파
	if isExist == false || isShorter == true {

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
			fmt.Println("--Sending Adv to Neighbor--")
			fmt.Println("From: ", uc.broker.Id)
			fmt.Println("To:   ", neighbor.Id)
			fmt.Println("-----------------------")
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

// go routine 2
// func (uc *brokerUsecase) DoSubscriptionQueue() {
// 	for {
// 		select {
// 		case reqPrtItem := <-uc.SubscriptionQueue:
// 			uc.SendSubscription(
// 				reqPrtItem,
// 			)
// 		}
// 	}
// }

// func (uc *brokerUsecase) PushSubscriptionToQueue(req *SubscriptionRequest) {
// 	uc.SubscriptionQueue <- req
// }

func (uc *brokerUsecase) SendSubscription(subReq *model.MessageRequest) error {
	// simulator.IncreaseCpuUsage()
	simulator.IncreaseMemoryUsage(100)

	newRequestPerformanceInfo := subReq.PerformanceInfo
	newRequestPerformanceInfo = append(newRequestPerformanceInfo, uc.GetPerformanceInfo())

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

	for _, item := range uc.broker.SRT {
		// advertisement에 subscription이 포함되는 경우:
		// PRT에 추가하고, 해당 advertisement의 last hop으로 subscription 전달함.
		if item.Advertisement.Subject == reqPrtItem.Subscription.Subject &&
			item.Advertisement.Operator == reqPrtItem.Subscription.Operator {
			advValue, _ := strconv.ParseFloat(item.Advertisement.Value, 64)
			subValue, _ := strconv.ParseFloat(reqPrtItem.Subscription.Value, 64)

			if (item.Advertisement.Operator == ">" && advValue >= subValue) ||
				(item.Advertisement.Operator == ">=" && advValue >= subValue) ||
				(item.Advertisement.Operator == "<" && advValue <= subValue) ||
				(item.Advertisement.Operator == "<=" && advValue <= subValue) {
				// PRT에 추가
				uc.broker.PRTmutex.Lock()
				uc.broker.PRT = append(uc.broker.PRT, reqPrtItem)
				uc.broker.PRTmutex.Unlock()

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
				for _, lastHop := range item.LastHop {
					fmt.Printf("%s\n", lastHop.Id)
				}

				for _, lastHop := range item.LastHop {
					fmt.Println("--Routing Sub via SRT--")
					fmt.Println("From: ", uc.broker.Id)
					fmt.Println("To:   ", lastHop.Id)
					fmt.Println("-----------------------")

					// 해당하는 advertisement를 보낸 publisher에게 도달한 경우: 전달 완료
					// if lastHop.NodeType == constants.PUBLISHER {
					if lastHop.Id == item.Identifier.SenderId {
						fmt.Printf("Subscription reached publisher %s\n", lastHop.Id)
						break
					}
					if subReq.Port == lastHop.Port {
						break
					}
					// 	for _, publisher := range uc.Publishers {
					// 		if publisher.ID == lastHop.ID {
					// 			// RPCNotifyPublisher
					// 		}
					// 	}

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
// func (uc *brokerUsecase) DoPublicationQueue() {
// 	for {
// 		select {
// 		case reqItem := <-uc.PublicationQueue:
// 			uc.SendPublication(
// 				reqItem,
// 			)
// 		}
// 	}
// }

// func (uc *brokerUsecase) PushPublicationToQueue(req *PublicationRequest) {
// 	uc.PublicationQueue <- req
// }

func (uc *brokerUsecase) SendPublication(pubReq *model.MessageRequest) error {
	// simulator.IncreaseCpuUsage()
	simulator.IncreaseMemoryUsage(100)

	newRequestPerformanceInfo := pubReq.PerformanceInfo
	newRequestPerformanceInfo = append(newRequestPerformanceInfo, uc.GetPerformanceInfo())

	uc.hopLogger.WithFields(logrus.Fields{
		"Node":            uc.broker.Id,
		"PerformanceInfo": newRequestPerformanceInfo,
	}).Infof("Publication(%s %s %s) %s", pubReq.Subject, pubReq.Operator, pubReq.Value, pubReq.MessageId)

	for _, item := range uc.broker.PRT {
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
				fmt.Println("From: ", uc.broker.Id)
				fmt.Println("To:   ", lastHop.Id)
				fmt.Println("-----------------------")

				// 해당하는 subscription을 보낸 subscriber에게 도달한 경우: 전달 완료
				// if lastHop.NodeType == constants.SUBSCRIBER {
				if lastHop.Id == item.Identifier.SenderId {
					fmt.Printf("Publication reached subscriber %s\n", lastHop.Id)

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
						log.Fatalf("error: %v", err)
						continue
					}
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
		interArrivalTime := broker.InterArrivalTime

		uc.tickLogger.WithFields(logrus.Fields{
			"CPU":                cpu,
			"Memory":             memory,
			"Queue Length":       queueLength,
			"Queue Time":         fmt.Sprintf("%f", queueTime.Seconds()*1000),   // 단위: ms, 소수점 아래 6자리
			"Service Time":       fmt.Sprintf("%f", serviceTime.Seconds()*1000), // 단위: ms, 소수점 아래 6자리
			"Response Time":      fmt.Sprintf("%f", responseTime.Seconds()*1000),
			"Throughput":         fmt.Sprintf("%.6f", throughput),
			"Inter-Arrival Time": fmt.Sprintf("%f", interArrivalTime.Seconds()*1000), // 단위: ms, 소수점 아래 6자리
			"Bottleneck":         uc.CheckBottleneck(&bottleneck, memory, responseTime),
		}).Info("Performance Metrics")

		LogBrokerInfoToLogfile(broker)
	}
}

type bottleneckStatus struct {
	start           time.Time
	preResponseTime time.Duration
	state           bool
}

func (uc *brokerUsecase) CheckBottleneck(bottleneckStatus *bottleneckStatus, memory uint64, responseTime time.Duration) bool {
	// 메모리 사용량이 100000000이면서 response time이 증가하는 상태가 n초 이상 지속되면 병목으로 판단
	var n time.Duration = 3
	if memory >= 100000000 && responseTime > bottleneckStatus.preResponseTime {
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
	bottleneckStatus.preResponseTime = responseTime
	return bottleneckStatus.state
}

func (uc *brokerUsecase) GetPerformanceInfo() *model.PerformanceInfo {
	broker := uc.broker
	cpu, mem := utils.Utilization()
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
		Memory:           fmt.Sprint(mem),
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

func LogBrokerInfoToLogfile(broker *model.Broker) {
	// save to log file (broker_id.log)
	// if file doesn't exist, create it, otherwise append to the file
	logFile, err := os.OpenFile("log/"+broker.Port+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	brokerInfoString := "==Broker Info==\n"
	brokerInfoString += fmt.Sprintf("Id: %s\n", broker.Id)
	brokerInfoString += fmt.Sprintf("Ip: %s\n", broker.Ip)
	brokerInfoString += fmt.Sprintf("Port: %s\n", broker.Port)

	publishersString := "==Publishers==\n"
	for _, publisher := range broker.Publishers {
		publishersString += fmt.Sprintf("Id: %s\n", publisher.Id)
		publishersString += fmt.Sprintf("Ip: %s\n", publisher.Ip)
		publishersString += fmt.Sprintf("Port: %s\n", publisher.Port)
	}
	subscribersString := "==Subscribers==\n"
	for _, subscriber := range broker.Subscribers {
		subscribersString += fmt.Sprintf("Id: %s\n", subscriber.Id)
		subscribersString += fmt.Sprintf("Ip: %s\n", subscriber.Ip)
		subscribersString += fmt.Sprintf("Port: %s\n", subscriber.Port)
	}
	brokerString := "==Brokers==\n"
	for _, neighbor := range broker.Brokers {
		brokerString += fmt.Sprintf("Id: %s\n", neighbor.Id)
		brokerString += fmt.Sprintf("Ip: %s\n", neighbor.Ip)
		brokerString += fmt.Sprintf("Port: %s\n", neighbor.Port)
	}
	itemsSrt := "==SRT==\n"
	for _, item := range broker.SRT {
		itemsSrt += fmt.Sprintf("Adv: %s %s %s (%s) | %s\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.Identifier.MessageId, item.Identifier.SenderId)
		for i := 0; i < len(item.LastHop); i++ {
			itemsSrt += fmt.Sprintf("%s | %s | %d \n", item.LastHop[i].Id, item.LastHop[i].NodeType, item.HopCount)
		}
		itemsSrt += fmt.Sprintf("----------------------------\n")
	}

	itemsPrt := "==PRT==\n"

	for _, item := range broker.PRT {
		itemsPrt += fmt.Sprintf("Sub: %s %s %s (%s) | %s\n", item.Subscription.Subject, item.Subscription.Operator, item.Subscription.Value, item.Identifier.MessageId, item.Identifier.SenderId)
		for i := 0; i < len(item.LastHop); i++ {
			itemsPrt += fmt.Sprintf("%s | %s\n", item.LastHop[i].Id, item.LastHop[i].NodeType)
		}
		itemsPrt += fmt.Sprintf("----------------------------\n")
	}

	log.Printf(brokerInfoString + publishersString + subscribersString + brokerString + itemsSrt + itemsPrt)
}
