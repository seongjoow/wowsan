package constants

import (
	"time"
)

var (
	// node type
	BROKER     = "broker"
	PUBLISHER  = "publisher"
	SUBSCRIBER = "subscriber"

	// message type
	ADVERTISEMENT = "advertisement"
	SUBSCRIPTION  = "subscription"
	PUBLICATION   = "publication"
)

var (
	DefaultSleepTime                      = 3 * time.Minute
	DefaultMessageServiceSleepTime        = 300 * time.Millisecond
	DefaultBroker3MessageServiceSleepTime = 100 * time.Millisecond
	MessageServiceSleepTime               = []time.Duration{5 * time.Second, 1 * time.Millisecond, 5 * time.Second, 1 * time.Millisecond}
	//when do message queue start time > DoMessageChangeSleepTime, service sleep time will change to next
	// (e.g. when do message queue start time > DoMessageChangeSleepTime[0],
	//service sleep time will change to MessageServiceSleepTime[0])
	DoMessageChangeSleepTime = []time.Duration{1 * time.Minute, 10 * time.Minute, 1 * time.Minute, 10 * time.Minute}
	// DoMessageChangeSleepTime = []time.Duration{10 * time.Second, 1 * time.Second, 10 * time.Second, 1 * time.Second}
	PublisherTimePauseTime = DoMessageChangeSleepTime
)
