package main

import (
	"fmt"
	"log"
	"net"
	"wowsan/pkg/broker"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/broker"

	grpcClient "wowsan/pkg/broker/transport"

	grpc "google.golang.org/grpc"

	_cli "wowsan/cmd/broker/seed/cli"
	logger "wowsan/pkg/logger"
)

var seedBrokers = []string{
	"50051",
	"50052",
	"50053",
}

var Brokers = []*model.Broker{}

func initSeed(port string) *model.Broker {
	lis, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()

	id := "localhost" + ":" + port
	l := logger.NewLogger(port)
	localBrokerModel := model.NewBroker(id, "localhost", port, l)
	server := broker.NewBrokerRPCServer(localBrokerModel)

	pb.RegisterBrokerServiceServer(s, server)
	go s.Serve(lis)
	go localBrokerModel.DoMessageQueue()
	// go localBrokerModel.DoAdvertisementQueue()
	// go localBrokerModel.DoSubscriptionQueue()
	// go localBrokerModel.DoPublicationQueue()

	fmt.Printf("Broker server listening at %v\n", lis.Addr())

	//add to Brokes
	return localBrokerModel
}

func main() {
	var isReady bool
	// var isReady = make([]bool, len(seedBrokers))
	var Brokers = []*model.Broker{}
	rpcBrokerClient := grpcClient.NewBrokerClient()

	for index, port := range seedBrokers {
		broker := initSeed(port)
		// isReady[index] = true
		if index == len(seedBrokers)-1 {
			isReady = true
			fmt.Printf("all broker is ready: %v\n", isReady)
		}
		Brokers = append(Brokers, broker)
	}

	go func() {
		if isReady {
			for index, broker := range Brokers {
				var addTo *model.Broker
				if index+1 >= len(Brokers) {
					addTo = Brokers[0]
				} else {
					addTo = Brokers[index+1]
				}
				rpcBrokerClient.RPCAddBroker(
					addTo.Ip,
					addTo.Port,
					broker.Id,
					broker.Ip,
					broker.Port,
				)
				broker.AddBroker(addTo.Id, addTo.Ip, addTo.Port)
			}
			// close the goroutine
			return
		}
	}()

	_cli.SeedCliLoop(rpcBrokerClient, Brokers)
}
