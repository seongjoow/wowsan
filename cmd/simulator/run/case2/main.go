package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"wowsan/constants"
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
		Duration:        120 * 60,
		Mu:              sync.Mutex{},
		Start:           time.Now(),
		AdvDuration:     120 * 60, // 1 hour
		AdvLambda:       0.1,      // 2 seconds
		PubLambda:       0.1,      // 2 seconds
		AdvControlChans: []chan string{},
		PubControlChan:  make(chan string, 10),
		NotStarted:      []chan string{},
		Running:         []chan string{},
		Paused:          []chan string{},
	}
}
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

func simulatorRandomPause() {
	// new pub server
	pubServers1 := NewPubServers()
	pubServers2 := NewPubServers()
	pubServers3 := NewPubServers()
	pubServers4 := NewPubServers()
	pubServers5 := NewPubServers()
	pubServerLoop := NewPubServers()

	pubServers5.AdvLambda = 4.5
	pubServers5.PubLambda = 2
	pubServerLoop.AdvLambda = 4.5
	pubServerLoop.PubLambda = 2

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

	publisherCount := 30
	for i := 0; i < publisherCount; i++ {
		startPubServers(&pubServerLoop, 60010, "50005")
		time.Sleep(1 * time.Second)
	}
	controlAllServers(&pubServerLoop, simulator.START)

	// for while to random pause, start, and resume until pubCloseTime
	pubCloseTime := time.Now().Add(150 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	// actions := []string{simulator.START, simulator.PAUSE, simulator.RESUME, simulator.PAUSE}
	actions := []string{simulator.PAUSE, simulator.RESUME}
	nextAction := 0
	level := 2
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// time.Sleep(1 * time.Minute)
	config, _ := constants.GetConfig("case2")
	time.Sleep(config.PublisherTimePauseTime[0] + config.PublisherTimePauseTime[1] + config.DefaultSleepTime)

	for range ticker.C {
		if time.Now().After(pubCloseTime) {
			break
		}
		action := actions[nextAction] // 0: pause, 1: resume
		if nextAction == len(actions)-1 {
			nextAction = 0
		} else {
			nextAction++
		}
		// 100%
		randPercentage := 1.0
		switch action {
		case simulator.PAUSE:
			fmt.Printf("Pausing %f%% of running servers...\n", randPercentage*100)
			controlRandomServers(&pubServerLoop, randPercentage, simulator.PAUSE)
		case simulator.RESUME:
			fmt.Printf("Resuming %f%% paused servers...\n", randPercentage*100)
			controlRandomServers(&pubServerLoop, randPercentage, simulator.RESUME)
		}
		time.Sleep(config.PublisherTimePauseTime[level])
		if level+1 >= len(config.PublisherTimePauseTime) {
			level = 0
		} else {
			level++
		}
	}
}

func main() {
	simulatorRandomPause()
}
