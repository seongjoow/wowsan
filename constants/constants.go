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

type MessageServiceConfig struct {
	DefaultSleepTime                      time.Duration
	DefaultMessageServiceSleepTime        time.Duration
	DefaultBroker3MessageServiceSleepTime time.Duration
	MessageServiceSleepTime               []time.Duration
	DoMessageChangeSleepTime              []time.Duration
	PublisherTimePauseTime                []time.Duration
}

var configs = map[string]MessageServiceConfig{
	"case1": {
		DefaultSleepTime:                      3 * time.Minute,
		DefaultMessageServiceSleepTime:        300 * time.Millisecond,
		DefaultBroker3MessageServiceSleepTime: 100 * time.Millisecond,
		MessageServiceSleepTime:               []time.Duration{5 * time.Second, 100 * time.Millisecond, 5 * time.Second, 300 * time.Millisecond, 5 * time.Second, 200 * time.Millisecond, 5 * time.Second, 300 * time.Millisecond, 5 * time.Second, 200 * time.Millisecond},
		DoMessageChangeSleepTime:              []time.Duration{1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 10 * time.Minute, 5 * time.Minute, 10 * time.Minute, 3 * time.Minute, 5 * time.Minute},
		PublisherTimePauseTime:                []time.Duration{1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 10 * time.Minute, 5 * time.Minute, 10 * time.Minute, 3 * time.Minute, 5 * time.Minute},
	},
	"case2": {
		DefaultSleepTime:                      3 * time.Minute,
		DefaultMessageServiceSleepTime:        300 * time.Millisecond,
		DefaultBroker3MessageServiceSleepTime: 100 * time.Millisecond,
		MessageServiceSleepTime:               []time.Duration{5 * time.Second, 100 * time.Millisecond, 5 * time.Second, 300 * time.Millisecond, 5 * time.Second, 200 * time.Millisecond, 5 * time.Second, 300 * time.Millisecond, 5 * time.Second, 200 * time.Millisecond},
		DoMessageChangeSleepTime:              []time.Duration{1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 10 * time.Minute, 5 * time.Minute, 10 * time.Minute, 3 * time.Minute, 5 * time.Minute},
		PublisherTimePauseTime:                []time.Duration{1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 10 * time.Minute, 5 * time.Minute, 10 * time.Minute, 3 * time.Minute, 5 * time.Minute},
	},
	"case3": {
		DefaultSleepTime:                      3 * time.Minute,
		DefaultMessageServiceSleepTime:        5 * time.Second,
		DefaultBroker3MessageServiceSleepTime: 100 * time.Millisecond,
		MessageServiceSleepTime:               []time.Duration{100 * time.Millisecond, 500 * time.Millisecond, 100 * time.Millisecond, 1000 * time.Millisecond, 100 * time.Millisecond, 1500 * time.Millisecond, 100 * time.Millisecond, 2000 * time.Millisecond, 100 * time.Millisecond, 2500 * time.Millisecond},
		DoMessageChangeSleepTime:              []time.Duration{1 * time.Minute, 15 * time.Minute, 2 * time.Minute, 10 * time.Minute, 3 * time.Minute, 15 * time.Minute, 5 * time.Minute, 10 * time.Minute, 3 * time.Minute, 5 * time.Minute},
		PublisherTimePauseTime:                []time.Duration{1 * time.Minute, 15 * time.Minute, 2 * time.Minute, 10 * time.Minute, 3 * time.Minute, 15 * time.Minute, 5 * time.Minute, 10 * time.Minute, 3 * time.Minute, 5 * time.Minute},
	},
	"case4": {},
}

func GetConfig(key string) (MessageServiceConfig, bool) {
	config, exists := configs[key]
	if !exists {
		config = configs["case1"]
	}
	return config, exists
}
