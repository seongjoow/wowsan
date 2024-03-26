package main

import (
	"flag"
	"wowsan/pkg/subscriber/cli"
	"wowsan/pkg/subscriber/service"
)

func main() {
	// ip := flag.String("ip", "", "ID of this node.")
	ip := "localhost"
	port := flag.String("port", "", "Port that this node should listen on.")
	flag.Parse()

	subscriberService := service.NewSubscriberService(ip, *port)

	// cli.ExecutionLoop(*ip, *port)
	cli.ExecutionLoop(subscriberService)
}
