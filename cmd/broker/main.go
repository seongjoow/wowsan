package main

import (
	"flag"
	cli "wowsan/pkg/broker/cli"
	"wowsan/pkg/broker/service"
)

func main() {
	// ip := flag.String("ip", "", "ID of this node.")
	ip := "localhost"
	port := flag.String("port", "", "Port that this node should listen on.")
	dirIndex := flag.String("dir_index", "", "Directory index for log files.")
	flag.Parse()

	brokerService := service.NewBrokerService(ip, *port, *dirIndex)

	// cli.ExecutionLoop(*ip, *port)
	cli.ExecutionLoop(brokerService)
}
