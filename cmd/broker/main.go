package main

import (
	"flag"
	"wowsan/pkg/broker/cli"
)

func main() {
	// ip := flag.String("ip", "", "ID of this node.")
	ip := "localhost"
	port := flag.String("port", "", "Port that this node should listen on.")
	flag.Parse()
	// cli.ExecutionLoop(*ip, *port)
	cli.ExecutionLoop(ip, *port)
}
