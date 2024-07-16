package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"wowsan/pkg/simulator"
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
	NotStarted      []chan string
	Running         []chan string
	Paused          []chan string
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
		NotStarted:      []chan string{},
		Running:         []chan string{},
		Paused:          []chan string{},
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

	var targetChans []chan string

	switch command {
	case simulator.PAUSE:
		targetChans = pubServers.Running
		pubServers.Paused = append(pubServers.Paused, targetChans...)
		pubServers.Running = nil

	case simulator.RESUME:
		targetChans = pubServers.Paused
		pubServers.Running = append(pubServers.Running, targetChans...)
		pubServers.Paused = nil

	case simulator.START:
		targetChans = pubServers.NotStarted
		pubServers.Running = append(pubServers.Running, targetChans...)
		pubServers.NotStarted = nil
	}

	for _, ch := range targetChans {
		ch <- command
	}
}

func controlRandomServers(pubServers *PubServers, percentage float64, command string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var targetChans []chan string

	switch command {
	case simulator.PAUSE:
		targetChans = pubServers.Running
	case simulator.RESUME:
		targetChans = pubServers.Paused
	case simulator.START:
		targetChans = pubServers.NotStarted
	}

	count := int(float64(len(targetChans)) * percentage)
	indices := r.Perm(len(targetChans))[:count]
	selectedChans := make([]chan string, count)

	for i, idx := range indices {
		selectedChans[i] = targetChans[idx]
		targetChans[idx] <- command
	}

	switch command {
	case simulator.PAUSE:
		pubServers.Paused = append(pubServers.Paused, selectedChans...)
		pubServers.Running = removeChans(pubServers.Running, selectedChans)

	case simulator.RESUME:
		pubServers.Running = append(pubServers.Running, selectedChans...)
		pubServers.Paused = removeChans(pubServers.Paused, selectedChans)

	case simulator.START:
		pubServers.Running = append(pubServers.Running, selectedChans...)
		pubServers.NotStarted = removeChans(pubServers.NotStarted, selectedChans)
	}
}

func removeChans(original, toRemove []chan string) []chan string {
	remaining := []chan string{}
	toRemoveMap := make(map[chan string]bool)
	for _, ch := range toRemove {
		toRemoveMap[ch] = true
	}
	for _, ch := range original {
		if !toRemoveMap[ch] {
			remaining = append(remaining, ch)
		}
	}
	return remaining
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
	pubServers.NotStarted = append(pubServers.NotStarted, controlChan)
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

// func simulator2minsPause2minsResume() {
// 	initPubDuration := 10 * time.Second

// 	pubServers1 := NewPubServers()
// 	pubServers2 := NewPubServers()
// 	pubServers3 := NewPubServers()
// 	pubServers4 := NewPubServers()
// 	pubServers5 := NewPubServers()
// 	// pubServers5.AdvLambda = 1.0

// 	pubServerLoop := NewPubServers()

// 	startPubServers(&pubServers1, 60001, "50001") //broker 1
// 	startPubServers(&pubServers2, 60002, "50001") //broker 1
// 	startPubServers(&pubServers3, 60003, "50003") //broker 3
// 	startPubServers(&pubServers4, 60004, "50004") //broker 4
// 	startPubServers(&pubServers5, 60005, "50005") //broker 5

// 	controlAllServers(&pubServers1, simulator.START)
// 	time.Sleep(initPunDuration)
// 	controlAllServers(&pubServers2, simulator.START)
// 	time.Sleep(initPunDuration)
// 	controlAllServers(&pubServers3, simulator.START)
// 	time.Sleep(initPunDuration)
// 	controlAllServers(&pubServers4, simulator.START)
// 	time.Sleep(initPunDuration)
// 	controlAllServers(&pubServers5, simulator.START)
// 	time.Sleep(initPunDuration)

// 	// for loop every 1 second
// 	checkDuration := 1 * time.Second
// 	ticker := time.NewTicker(checkDuration)

// 	startTime := time.Now()
// 	simulateTimer := time.Now() // to calculate whether or not to pause the simulation, will reset when the pause time is reached
// 	pubCloseTime := startTime.Add(120 * time.Minute)
// 	addNewPubServerTime := startTime.Add(initPunDuration) // to calculate whether or not to start new pub server

// 	for range ticker.C {
// 		simulateEndTime := simulateTimer.Add(2 * time.Minute)
// 		pauseEndTime := simulateEndTime.Add(2 * time.Minute)

// 		isPaused := false

// 		if time.Now().Before(simulateEndTime) { // new pub server and send message
// 			if time.Now().After(addNewPubServerTime) {
// 				startPubServers(&pubServerLoop, 60010, "50003")
// 				controlAllServers(&pubServerLoop, simulator.START)
// 				addNewPubServerTime = time.Now().Add(initPunDuration)
// 			}

// 		} else if time.Now().Before(pauseEndTime) { // pause
// 			if !isPaused {
// 				fmt.Printf("Pausing broker 3 message generation...\n")
// 				controlAllServers(&pubServers3, simulator.PAUSE)
// 				controlAllServers(&pubServerLoop, simulator.PAUSE)
// 				isPaused = true
// 			}
// 		} else {
// 			fmt.Printf("Resuming broker 3 message generation...\n")
// 			// reset start time
// 			controlAllServers(&pubServers3, simulator.RESUME)
// 			controlAllServers(&pubServerLoop, simulator.RESUME)

// 			isPaused = false
// 			simulateTimer = pauseEndTime
// 		}

// 		if time.Now().After(pubCloseTime) {
// 			break
// 		}
// 	}
// }

func simulatorRandomPause() {
	// new pub server
	pubServers1 := NewPubServers()
	pubServers2 := NewPubServers()
	pubServers3 := NewPubServers()
	pubServers4 := NewPubServers()
	pubServers5 := NewPubServers()
	pubServerLoop := NewPubServers()

	startPubServers(&pubServers1, 60001, "50001") //broker 1
	startPubServers(&pubServers2, 60002, "50001") //broker 1
	startPubServers(&pubServers3, 60003, "50003") //broker 3
	startPubServers(&pubServers4, 60004, "50004") //broker 4
	startPubServers(&pubServers5, 60005, "50005") //broker 5

	controlAllServers(&pubServers1, simulator.START)
	controlAllServers(&pubServers2, simulator.START)
	controlAllServers(&pubServers3, simulator.START)
	controlAllServers(&pubServers4, simulator.START)
	controlAllServers(&pubServers5, simulator.START)

	publisherCount := 3
	for i := 0; i < publisherCount; i++ {
		startPubServers(&pubServerLoop, 60010, "50003")
	}
	controlAllServers(&pubServerLoop, simulator.START)

	// for while to random pause, start, and resume until pubCloseTime
	pubCloseTime := time.Now().Add(120 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	// actions := []string{simulator.START, simulator.PAUSE, simulator.RESUME, simulator.PAUSE}
	actions := []string{simulator.PAUSE, simulator.RESUME}
	nextAction := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	init := false

	for range ticker.C {
		pauseRandomInterval := time.Duration(r.Intn(1)+3) * time.Minute
		resumeRandomInterval := time.Duration(r.Intn(1)+1) * time.Minute
		if !init {
			// for the first time
			init = true
			time.Sleep(resumeRandomInterval)
			continue
		}
		// randomInterval := time.Duration(r.Intn(4)+1) * time.Minute
		if time.Now().After(pubCloseTime) {
			break
		}
		// action := actions[nextAction] // 0: start, 1: pause, 2: resume, 3: pause
		action := actions[nextAction] // 0: pause, 1: resume
		if nextAction == len(actions)-1 {
			nextAction = 0
		} else {
			nextAction++
		}
		// randPercentage := 0.5 + r.Float64()*0.5 // 0.3 + [0, 0.4) -> [0.3, 0.7)
		//70~100%
		randPercentage := 0.7 + r.Float64()*0.2
		// 100%
		// randPercentage := 1.0
		switch action {
		case simulator.PAUSE:
			fmt.Printf("Pausing %f%% of running servers...\n", randPercentage*100)
			controlRandomServers(&pubServerLoop, randPercentage, simulator.PAUSE)
			controlRandomServers(&pubServers1, randPercentage, simulator.PAUSE)
			controlRandomServers(&pubServers2, randPercentage, simulator.PAUSE)
			controlRandomServers(&pubServers3, randPercentage, simulator.PAUSE)
			controlRandomServers(&pubServers4, randPercentage, simulator.PAUSE)
			controlRandomServers(&pubServers5, randPercentage, simulator.PAUSE)
			time.Sleep(pauseRandomInterval)
		case simulator.RESUME:
			fmt.Printf("Resuming %f%% paused servers...\n", randPercentage*100)
			controlRandomServers(&pubServerLoop, randPercentage, simulator.RESUME)
			controlRandomServers(&pubServers1, randPercentage, simulator.RESUME)
			controlRandomServers(&pubServers2, randPercentage, simulator.RESUME)
			controlRandomServers(&pubServers3, randPercentage, simulator.RESUME)
			controlRandomServers(&pubServers4, randPercentage, simulator.RESUME)
			controlRandomServers(&pubServers5, randPercentage, simulator.RESUME)
			time.Sleep(resumeRandomInterval)
			// case simulator.START:
			// 	if &pubServerLoop.NotStarted == nil {
			// 		fmt.Printf("Resuming %f%% paused servers...\n", randPercentage*100)
			// 		controlRandomServers(&pubServerLoop, randPercentage, simulator.RESUME)
			// 	} else {
			// 		fmt.Printf("Starting %f%% not started servers...\n", randPercentage*100)
			// 		controlRandomServers(&pubServerLoop, randPercentage, simulator.START)
			// 	}
		}
	}
}

func main() {
	// simulator2minsPause2minsResume()
	simulatorRandomPause()
}
