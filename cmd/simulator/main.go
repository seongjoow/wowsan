package main

import (
	"math/rand"
	"net/http"
	"time"
	"wowsan/pkg/simulator"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	advLambda := float64(1) / 8 // 단위 시간당 평균 호출 횟수 λ (예: λ = 0.333 이면 3초에 한 번 호출)
	subLambda := float64(1) / 5
	pubLambda := float64(1) / 3
	duration := 30 // 시뮬레이션 할 총 시간(초)

	go simulator.RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50001", "localhost", "60001")
	go simulator.RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50002", "localhost", "60002")
	go simulator.RunPublisherSimulation(duration, advLambda, pubLambda, "localhost", "50003", "localhost", "60003")
	// time.Sleep(3 * time.Second)
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50004", "localhost", "60004")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50005", "localhost", "60005")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50006", "localhost", "60006")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50007", "localhost", "60007")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50008", "localhost", "60008")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50009", "localhost", "60009")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50010", "localhost", "60010")

	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50004", "localhost", "60014")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50005", "localhost", "60015")
	go simulator.RunSubscriberSimulation(duration, subLambda, "localhost", "50006", "localhost", "60016")

	// wait until all goroutines are finished
	time.Sleep(time.Duration(duration) * time.Second)
	http.PostForm("http://localhost:8080/simulator/done", map[string][]string{})
}
