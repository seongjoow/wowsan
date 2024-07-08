package main

import (
	"fmt"
	"sync"
	"time"
	"wowsan/pkg/simulator" // 确保导入路径正确
)

type PubServers struct {
	Servers     []string
	Duration    int
	Mu          sync.Mutex
	Start       time.Time
	AdvDuration int
	AdvLambda   float64
	PubLambda   float64
	// AdvControlChan chan string
	AdvControlChans []chan string
	PubControlChan  chan string
}

func NewPubServers() PubServers {
	return PubServers{
		Servers:         []string{},
		Duration:        60 * 60,
		Mu:              sync.Mutex{},
		Start:           time.Now(),
		AdvDuration:     60 * 60, // 1 hour
		AdvLambda:       2.0,     // 2 seconds
		PubLambda:       2.0,     // 2 seconds
		AdvControlChans: []chan string{},
		PubControlChan:  make(chan string, 10),
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

func controlAllServers(pubServers *PubServers, command string) {
	pubServers.Mu.Lock()
	defer pubServers.Mu.Unlock()

	for _, ch := range pubServers.AdvControlChans {
		ch <- command
	}
}

func startPubServers(pubServers *PubServers, pubStartPort int, brokerPort string) {
	var port int
	if len(pubServers.Servers) > 0 {
		port = pubStartPort + len(pubServers.Servers) + 1
	} else {
		port = pubStartPort
	}

	controlChan := make(chan string, 10)
	pubServers.Mu.Lock()
	pubServers.AdvControlChans = append(pubServers.AdvControlChans, controlChan)
	pubServers.Mu.Unlock()

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
		controlChan,
		pubServers.PubControlChan,
	)

	serverAddress := fmt.Sprintf(":%d", port)
	pubServers.Mu.Lock()
	pubServers.Servers = append(pubServers.Servers, serverAddress)
	pubServers.Mu.Unlock()
	fmt.Printf("Server started at %s\n", serverAddress)
}

func main() {
	initPunDuration := 10 * time.Second

	pubServers1 := NewPubServers()
	pubServers2 := NewPubServers()
	pubServers3 := NewPubServers()
	pubServers4 := NewPubServers()
	pubServers5 := NewPubServers()
	// pubServers5.AdvLambda = 1.0

	pubServerLoop := NewPubServers()

	startPubServers(&pubServers1, 60001, "50001") //broker 1
	time.Sleep(initPunDuration)
	startPubServers(&pubServers2, 60002, "50001") //broker 1
	time.Sleep(initPunDuration)
	startPubServers(&pubServers3, 60003, "50003") //broker 3
	time.Sleep(initPunDuration)
	startPubServers(&pubServers4, 60004, "50004") //broker 4
	time.Sleep(initPunDuration)
	startPubServers(&pubServers5, 60005, "50005") //broker 5
	time.Sleep(initPunDuration)

	// for loop every 1 second
	checkDuration := 1 * time.Second
	ticker := time.NewTicker(checkDuration)

	startTime := time.Now()
	simulateTimer := time.Now() // to calculate whether or not to pause the simulation, will reset when the pause time is reached
	pubCloseTime := startTime.Add(120 * time.Minute)
	addNewPubServerTime := startTime.Add(initPunDuration) // to calculate whether or not to start new pub server

	for range ticker.C {
		simulateEndTime := simulateTimer.Add(2 * time.Minute)
		pauseEndTime := simulateEndTime.Add(2 * time.Minute)

		isPaused := false

		if time.Now().Before(simulateEndTime) { // new pub server and send message
			if time.Now().After(addNewPubServerTime) {
				startPubServers(&pubServerLoop, 60010, "50003")
				addNewPubServerTime = time.Now().Add(initPunDuration)
			}

		} else if time.Now().Before(pauseEndTime) { // pause
			if !isPaused {
				fmt.Printf("Pausing broker 3 message generation...\n")
				controlAllServers(&pubServers3, simulator.PAUSE)
				controlAllServers(&pubServerLoop, simulator.PAUSE)
				isPaused = true
			}
		} else {
			fmt.Printf("Resuming broker 3 message generation...\n")
			// reset start time
			controlAllServers(&pubServers3, simulator.RESUME)
			controlAllServers(&pubServerLoop, simulator.RESUME)

			isPaused = false
			simulateTimer = pauseEndTime
		}

		// exit
		if time.Now().After(pubCloseTime) {
			break
		}
	}
}
