package main

import (
	"math/rand"
	"time"
	publisher "wowsan/pkg/publisher/service"
	"wowsan/pkg/simulator"
	subscriber "wowsan/pkg/subscriber/service"
)

var subjectList = []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}

// getExpInterval 함수는 지수 분포를 사용하여 다음 호출까지의 대기 시간을 반환함
func getExpInterval(lambda float64) time.Duration {
	expRandom := rand.ExpFloat64() / lambda
	return time.Duration(expRandom * float64(time.Second))
}

// runSimulation 함수는 주어진 시간 동안 시뮬레이션을 실행함
func RunPublisherSimulation(durationSeconds int, advLambda float64, pubLambda float64, brokerIp string, brokerPort string, publisherId string, publisherIp string, publisherPort string) {
	publisherService := publisher.NewPublisherService(publisherIp, publisherPort)

	start := time.Now()
	end := start.Add(time.Duration(durationSeconds) * time.Second)

	subject := ""
	operator := ""
	value := ""
	selectedSubjectList := []string{}

	// Advertise 시뮬레이션 루프
	for time.Now().Before(end) {
		subject, operator, value, selectedSubjectList = simulator.AdvPredicateGenerator(subjectList)
		// value to string
		// strValue := fmt.Sprintf("%d", value)

		interval := getExpInterval(advLambda)
		time.Sleep(interval)

		// // 메세지 전달 함수 호출
		publisherService.PublisherUsecase.Adv(subject, operator, value, brokerIp, brokerPort)

		// brokerClient.RPCSendAdvertisement(
		// 	publisherModel.Broker.Ip,
		// 	publisherModel.Broker.Port,
		// 	publisherModel.Id,
		// 	publisherModel.Ip,
		// 	publisherModel.Port,
		// 	subject,
		// 	operator,
		// 	value,
		// 	constants.PUBLISHER,
		// 	hopCount,
		// 	uuid.NewV4().String(),
		// 	publisherId,
		// )
	}

	// Publisher 시뮬레이션 루프
	for time.Now().Before(end) {
		subject, operator, value := simulator.PubPredicateGenerator(selectedSubjectList)
		// strValue := fmt.Sprintf("%d", value) // value to string

		interval := getExpInterval(pubLambda)
		time.Sleep(interval)

		// 메세지 전달 함수 호출
		publisherService.PublisherUsecase.Pub(subject, operator, value, brokerIp, brokerPort)

		// brokerClient.RPCSendPublication(
		// 	publisherModel.Broker.Ip,
		// 	publisherModel.Broker.Port,
		// 	publisherModel.Id,
		// 	publisherModel.Ip,
		// 	publisherModel.Port,
		// 	subject,
		// 	operator,
		// 	value,
		// 	constants.PUBLISHER,
		// 	uuid.NewV4().String(),
		// )
	}
}

func RunSubscriberSimulation(durationSeconds int, lambda float64, brokerIp string, brokerPort string, subscriberId string, subscriberIp string, subscriberPort string) {
	subscriberService := subscriber.NewSubscriberService(subscriberIp, subscriberPort)

	start := time.Now()
	end := start.Add(time.Duration(durationSeconds) * time.Second)

	// Subscribe 시뮬레이션 루프
	for time.Now().Before(end) {
		subject, operator, value := simulator.SubPredicateGenerator()
		// strValue := fmt.Sprintf("%d", value) // value to string

		interval := getExpInterval(lambda)
		time.Sleep(interval)

		// 메세지 전달 함수 호출
		subscriberService.SubscriberUsecase.Sub(subject, operator, value, brokerIp, brokerPort)

		// brokerClient.RPCSendSubscription(
		// 	subscriberModel.Broker.Ip,
		// 	subscriberModel.Broker.Port,
		// 	subscriberModel.Id,
		// 	subscriberModel.Ip,
		// 	subscriberModel.Port,
		// 	subject,
		// 	operator,
		// 	value,
		// 	constants.SUBSCRIBER,
		// 	uuid.NewV4().String(),
		// 	subscriberId,
		// )
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	advLambda := float64(1) / 5 // 단위 시간당 평균 호출 횟수 λ (예: λ = 0.333 이면 3초에 한 번 호출)
	pubLambda := float64(1) / 2
	subLambda := float64(1) / 3
	duration := 20 // 시뮬레이션 할 총 시간(초)

	go RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50051", "id1", "localhost", "2221")
	// time.Sleep(3 * time.Second)
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50054", "id2", "localhost", "2222")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50056", "id3", "localhost", "2223")

	// block thread
	select {}
}
