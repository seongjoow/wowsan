package main

import (
	"flag"
	"wowsan/pkg/publisher/cli"
	"wowsan/pkg/publisher/service"
)

func main() {
	// ip := flag.String("ip", "", "ID of this node.")
	ip := "localhost"
	port := flag.String("port", "", "Port that this node should listen on.")
	flag.Parse()

	publisherService := service.NewPublisherService(ip, *port)

	// cli.ExecutionLoop(*ip, *port)
	cli.ExecutionLoop(publisherService)
}
