package main

import (
	"fmt"
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
func RunPublisherSimulation(durationSeconds int, advLambda float64, pubLambda float64, brokerIp string, brokerPort string, publisherIp string, publisherPort string) {
	publisherService := publisher.NewPublisherService(publisherIp, publisherPort)

	start := time.Now()
	end := start.Add(time.Duration(durationSeconds) * time.Second)

	// subject := ""
	// operator := ""
	// value := ""
	// selectedSubjectList := []string{}
	// 채널 생성
	subjectListChannel := make(chan []string)

	// Advertise 시뮬레이션 루프
	go func() {
		for time.Now().Before(end) {
			subject, operator, value, newSelectedSubjectList := simulator.AdvPredicateGenerator(subjectList)

			interval := getExpInterval(advLambda)
			time.Sleep(interval)

			// 채널을 통해 리스트 전송
			subjectListChannel <- newSelectedSubjectList

			// // 메세지 전달 함수 호출
			publisherService.PublisherUsecase.Adv(subject, operator, value, brokerIp, brokerPort)
		}

		close(subjectListChannel) // 고루틴 종료 시 채널 닫기
	}()

	// Publish 시뮬레이션 루프
	for time.Now().Before(end) {
		// 채널에서 데이터 수신
		selectedSubjectList, ok := <-subjectListChannel
		if !ok {
			break // 채널이 닫혔다면 루프 종료
		}

		fmt.Println(len(selectedSubjectList))

		if len(selectedSubjectList) == 0 {
			continue
		}

		subject, operator, value := simulator.PubPredicateGenerator(selectedSubjectList)

		interval := getExpInterval(pubLambda)
		time.Sleep(interval)

		// 메세지 전달 함수 호출
		publisherService.PublisherUsecase.Pub(subject, operator, value, brokerIp, brokerPort)
	}
}

func RunSubscriberSimulation(durationSeconds int, lambda float64, brokerIp string, brokerPort string, subscriberIp string, subscriberPort string) {
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
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	advLambda := float64(1) / 8 // 단위 시간당 평균 호출 횟수 λ (예: λ = 0.333 이면 3초에 한 번 호출)
	subLambda := float64(1) / 5
	pubLambda := float64(1) / 3
	duration := 30 // 시뮬레이션 할 총 시간(초)

	go RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50001", "localhost", "60001")
	go RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50002", "localhost", "60002")
	go RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50003", "localhost", "60003")
	// time.Sleep(3 * time.Second)
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50004", "localhost", "60004")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50005", "localhost", "60005")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50006", "localhost", "60006")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50007", "localhost", "60007")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50008", "localhost", "60008")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50009", "localhost", "60009")
	go RunSubscriberSimulation(duration, subLambda, "localhost", "50010", "localhost", "60010")

	// block thread
	select {}
}
