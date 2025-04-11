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
	"case1": { // 인접노드에서 drop
		DefaultSleepTime:                      3 * time.Minute,
		DefaultMessageServiceSleepTime:        300 * time.Millisecond,
		DefaultBroker3MessageServiceSleepTime: 100 * time.Millisecond,
		// stable, unstable, stable, unstable, ...
		MessageServiceSleepTime:  []time.Duration{200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond},
		DoMessageChangeSleepTime: []time.Duration{5 * time.Minute, 1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 4 * time.Minute, 5 * time.Minute, 3 * time.Minute},
		PublisherTimePauseTime:   []time.Duration{5 * time.Minute, 1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 4 * time.Minute, 5 * time.Minute, 3 * time.Minute},
	},
	"case2": { // 비인접노드에서 drop
		DefaultSleepTime:                      3 * time.Minute,
		DefaultMessageServiceSleepTime:        300 * time.Millisecond,
		DefaultBroker3MessageServiceSleepTime: 100 * time.Millisecond,
		// stable, unstable, stable, unstable, ...
		MessageServiceSleepTime:  []time.Duration{200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond, 200 * time.Millisecond, 500 * time.Millisecond},
		DoMessageChangeSleepTime: []time.Duration{5 * time.Minute, 4 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 4 * time.Minute, 5 * time.Minute, 3 * time.Minute},
		PublisherTimePauseTime:   []time.Duration{5 * time.Minute, 4 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 4 * time.Minute, 5 * time.Minute, 3 * time.Minute},
	},

	"case3": { // 인접노드에서 spike
		DefaultSleepTime:                      3 * time.Minute,
		DefaultMessageServiceSleepTime:        300 * time.Millisecond,
		DefaultBroker3MessageServiceSleepTime: 100 * time.Millisecond,
		// stable, unstable, stable, unstable, ...
		MessageServiceSleepTime:  []time.Duration{200 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond, 100 * time.Millisecond},
		DoMessageChangeSleepTime: []time.Duration{5 * time.Minute, 1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 4 * time.Minute},
		PublisherTimePauseTime:   []time.Duration{5 * time.Minute, 1 * time.Minute, 5 * time.Minute, 2 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 3 * time.Minute, 5 * time.Minute, 4 * time.Minute},
	},
	// "case4": { //비인접노드에서 drop
	// },
}

func GetConfig(key string) (MessageServiceConfig, bool) {
	config, exists := configs[key]
	if !exists {
		config = configs["case3"]
	}
	return config, exists
}
