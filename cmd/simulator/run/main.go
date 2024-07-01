package main

import (
	"fmt"
	"sync"
	"time"
	"wowsan/pkg/simulator" // 确保导入路径正确
)

type PubServers struct {
	Servers        []string
	Duration       int
	Mu             sync.Mutex
	Start          time.Time
	AdvDuration    int
	AdvLambda      float64
	PubLambda      float64
	AdvControlChan chan string
	PubControlChan chan string
}

func NewPubServers() PubServers {
	return PubServers{
		Servers:        []string{},
		Duration:       60 * 60,
		Mu:             sync.Mutex{},
		Start:          time.Now(),
		AdvDuration:    60 * 60, // 1 hour
		AdvLambda:      2.0,     // 2 seconds
		PubLambda:      2.0,     // 2 seconds
		AdvControlChan: make(chan string, 10),
		PubControlChan: make(chan string, 10),
	}
}

// func NewSubServers() SubServers {
// 	return SubServers{
// 		Servers:  []string{},
// 		Duration: 30 * 1,
// 		Mu:       sync.Mutex{},
// 		Start:    time.Now(),
// 	}
// }

// type SubServers struct {
// 	Servers  []string
// 	Duration time.Duration
// 	Mu       sync.Mutex
// 	Start    time.Time
// }

func startPubServers(pubServers *PubServers, pubStartPort int, brokerPort string) {
	port := pubStartPort + len(pubServers.Servers) + 1
	go simulator.RunPublisherSimulation(
		pubServers.AdvDuration,
		pubServers.Duration,
		pubServers.AdvLambda,
		pubServers.PubLambda,
		"localhost",
		brokerPort,
		"localhost",
		fmt.Sprintf("%d", port),
		[]string{"apple"},
		pubServers.AdvControlChan,
		pubServers.PubControlChan,
	)

	serverAddress := fmt.Sprintf(":%d", port)
	pubServers.Mu.Lock()
	pubServers.Servers = append(pubServers.Servers, serverAddress)
	pubServers.Mu.Unlock()
	fmt.Printf("Server started at %s\n", serverAddress)
}

func main() {
	sleepTime := 10 * time.Second
	startTime := time.Now()
	closeTime := startTime.Add(120 * time.Minute)

	ticker := time.NewTicker(sleepTime)

	pubServers1 := NewPubServers()
	pubServers2 := NewPubServers()
	pubServers3 := NewPubServers()
	pubServers4 := NewPubServers()
	pubServers5 := NewPubServers()
	// pubServers5.AdvLambda = 1.0

	pubServerLoop := NewPubServers()

	startPubServers(&pubServers1, 60001, "50001") //broker 1
	time.Sleep(sleepTime)
	startPubServers(&pubServers2, 60002, "50001") //broker 1
	time.Sleep(sleepTime)
	startPubServers(&pubServers3, 60003, "50003") //broker 3
	time.Sleep(sleepTime)
	startPubServers(&pubServers4, 60004, "50004") //broker 4
	time.Sleep(sleepTime)
	startPubServers(&pubServers5, 60005, "50005") //broker 5
	time.Sleep(sleepTime)

	// if start time is 10 minutes, set the pubServer3 pause time to 10 minutes, then resume it after 10 minutes.
	go func() {
		time.Sleep(10 * time.Minute)
		fmt.Print("========================")
		fmt.Println("pause")
		pubServers3.Mu.Lock()
		pubServers3.AdvControlChan <- simulator.PAUSE
		pubServers3.Mu.Unlock()
		fmt.Print("========================")
		time.Sleep(5 * time.Minute)
		fmt.Print("========================")
		fmt.Println("resume")
		pubServers3.Mu.Lock()
		pubServers3.AdvControlChan <- simulator.RESUME
		pubServers3.Mu.Unlock()
		fmt.Print("========================")
	}()

	for range ticker.C {
		startPubServers(&pubServerLoop, 60010, "50003")
		if time.Now().After(closeTime) {
			break
		}
	}
}
