package main

import (
	"math/rand"
	"net/http"
	"time"
	"wowsan/pkg/simulator"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	advLambda := float64(1) / (60 * 5) // 단위 시간당 평균 호출 횟수 λ (예: λ = 1/3 이면 3초에 한 번 호출)
	subLambda := float64(1) / (60 * 2)
	pubLambda := float64(1) / 30
	duration := 60 * 60    // 시뮬레이션 할 총 시간(초)
	advDuration := 60 * 30 // Advertise 시뮬레이션 시간(초)
	// pubDuration := 60 * 60 // Publish 시뮬레이션 시간(초)

	var subjectList1 = []string{"apple", "tesla", "microsoft", "amazon", "nvidia"}
	var subjectList2 = []string{"cj", "samsung", "hyundai", "lg", "sk"}
	var subjectList3 = []string{"meditox", "celltrion", "samsungbiologics", "lgchem", "hanmi"}
	var subjectList4 = []string{"s&p500", "nasdaq", "kospi", "kosdaq", "dowjones"}
	var subjectList5 = []string{"yg", "jyp", "sm", "cube", "hybe"}

	go simulator.RunPublisherSimulation(advDuration, duration, advLambda, pubLambda, "localhost", "50001", "localhost", "60001", subjectList1)
	go simulator.RunPublisherSimulation(advDuration, duration, advLambda, pubLambda, "localhost", "50002", "localhost", "60002", subjectList2)

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50003", "localhost", "60013", subjectList1)
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50004", "localhost", "60014", subjectList1)
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50005", "localhost", "60015", subjectList2)
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50006", "localhost", "60016", subjectList2)

	time.Sleep(60 * 10 * time.Second)
	// 10분 경과

	go simulator.RunPublisherSimulation(advDuration, duration, advLambda, pubLambda, "localhost", "50003", "localhost", "60003", subjectList3)

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50007", "localhost", "60007", subjectList3)

	time.Sleep(60 * 5 * time.Second)
	// 15분 경과

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50008", "localhost", "60008", subjectList3)
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50009", "localhost", "60009", subjectList3)

	time.Sleep(60 * 10 * time.Second)
	// 25분 경과

	go simulator.RunPublisherSimulation(advDuration, duration, advLambda, pubLambda, "localhost", "50004", "localhost", "60004", subjectList4)

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50010", "localhost", "60010", subjectList4)

	time.Sleep(60 * 10 * time.Second)
	// 35분 경과

	go simulator.RunPublisherSimulation(advDuration, duration, advLambda, pubLambda, "localhost", "50005", "localhost", "60005", subjectList5)

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50001", "localhost", "60011", subjectList5)

	time.Sleep(60 * 5 * time.Second)
	// 40분 경과

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50002", "localhost", "60012", subjectList5)

	time.Sleep(60 * 10 * time.Second)
	// 50분 경과
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50003", "localhost", "60013", subjectList1)
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50004", "localhost", "60014", subjectList2)

	// Wait until all goroutines are finished
	time.Sleep(time.Duration(duration) * time.Second)
	http.PostForm("http://localhost:8080/simulator/done", map[string][]string{})
}
