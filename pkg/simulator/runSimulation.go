package simulator

import (
	"fmt"
	//"math/rand"
	"time"

	publisher "wowsan/pkg/publisher/service"
	subscriber "wowsan/pkg/subscriber/service"
)

const (
	PAUSE  = "pause"
	RESUME = "resume"
)

// getExpInterval 함수는 지수 분포를 사용하여 다음 호출까지의 대기 시간을 반환함
func getExpInterval(lambda float64) time.Duration {
	// expRandom := rand.ExpFloat64() / lambda
	// return time.Duration(expRandom * float64(time.Second))

	// format the lambda to seconds
	return time.Duration(lambda * float64(time.Second))
}

// runSimulation 함수는 주어진 시간 동안 시뮬레이션을 실행함
func RunPublisherSimulation(advDurationSeconds, pubDurationSeconds int, advLambda, pubLambda float64, brokerIp, brokerPort, publisherIp, publisherPort string, subjectList []string, advControlChan, pubControlChan chan string) {
	publisherService := publisher.NewPublisherService(publisherIp, publisherPort)

	start := time.Now()
	advEnd := start.Add(time.Duration(advDurationSeconds) * time.Second)
	pubEnd := start.Add(time.Duration(pubDurationSeconds) * time.Second)

	// subject := ""
	// operator := ""
	// value := ""
	// selectedSubjectList := []string{}
	// 채널 생성
	subjectListChannel := make(chan []string)
	// Advertise 시뮬레이션 루프
	go func() {
		for time.Now().Before(advEnd) {
			// fmt.Printf("len of controlChan:%d [from %s]\n", len(advControlChan), publisherPort)
			select {
			case cmd := <-advControlChan:
				fmt.Printf("Advertisement received command: %s\n", cmd)
				if cmd == PAUSE {
					for time.Now().Before(advEnd) {
						fmt.Printf("Advertisement loop paused, port:%s\n", publisherPort)
						if RESUME == <-advControlChan {
							fmt.Printf("Advertisement loop resumed\n")
							break
						}
					}
				}
			case <-time.After(getExpInterval(advLambda)):

				fmt.Printf("Publishing...[from %s| to %s]\n", publisherPort, brokerPort)
				subject, operator, value, newSelectedSubjectList := AdvPredicateGenerator(subjectList)

				// 채널을 통해 리스트 전송
				subjectListChannel <- newSelectedSubjectList

				// // 메세지 전달 함수 호출
				publisherService.PublisherUsecase.Adv(subject, operator, value, brokerIp, brokerPort)
			}
		}

		close(subjectListChannel) // 고루틴 종료 시 채널 닫기
	}()

	// Publish 시뮬레이션 루프
	for time.Now().Before(pubEnd) {
		select {
		case cmd := <-pubControlChan:
			if cmd == PAUSE {
				for time.Now().Before(pubEnd) {
					fmt.Printf("Publication loop paused, port:%s\n", publisherPort)
					if RESUME == <-pubControlChan {
						fmt.Printf("Publication loop resumed\n")
						break
					}
				}
			}
		default:

			// 채널에서 데이터 수신
			selectedSubjectList, ok := <-subjectListChannel
			if !ok {
				break // 채널이 닫혔다면 루프 종료
			}

			fmt.Println(len(selectedSubjectList))

			if len(selectedSubjectList) == 0 {
				continue
			}

			subject, operator, value := PubPredicateGenerator(selectedSubjectList)

			interval := getExpInterval(pubLambda)
			time.Sleep(interval)

			// 메세지 전달 함수 호출
			publisherService.PublisherUsecase.Pub(subject, operator, value, brokerIp, brokerPort)
		}
	}
}

func RunSubscriberSimulation(durationSeconds int, lambda float64, brokerIp, brokerPort, subscriberIp, subscriberPort string, subjectList []string) {
	subscriberService := subscriber.NewSubscriberService(subscriberIp, subscriberPort)

	start := time.Now()
	end := start.Add(time.Duration(durationSeconds) * time.Second)

	// Subscribe 시뮬레이션 루프
	for time.Now().Before(end) {
		subject, operator, value := SubPredicateGenerator(subjectList)
		// strValue := fmt.Sprintf("%d", value) // value to string

		interval := getExpInterval(lambda)
		time.Sleep(interval)

		// 메세지 전달 함수 호출
		subscriberService.SubscriberUsecase.Sub(subject, operator, value, brokerIp, brokerPort)
	}
}
