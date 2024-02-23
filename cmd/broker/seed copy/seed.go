package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"wowsan/pkg/broker"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/broker"

	grpcClient "wowsan/pkg/broker/transport"

	grpc "google.golang.org/grpc"

	_cli "wowsan/cmd/broker/seed/cli"
	logger "wowsan/pkg/logger"
)

var Brokers = []*model.Broker{}

func brokerPortsGenerator(counts int) (ports []string) {
	for i := 0; i < counts; i++ {
		ports = append(ports, fmt.Sprintf("%d", 50051+i))
	}
	return
}

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

	seedBrokers := brokerPortsGenerator(10)
	MAX_NEIGHBOR := 4
	// AVG_NEIGHBOR := 0
	MIN_NEIGHBOR := 2

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
			for _, broker := range Brokers {
				// fmt.Printf("broker %v has %v neighbor\n", broker.Id, len(broker.Brokers))
				if len(broker.Brokers) > MAX_NEIGHBOR {
					// fmt.Printf("broker %v has enough neighbor\n", broker.Id)
					break
				} else if len(broker.Brokers) < MIN_NEIGHBOR {
					var addTo *model.Broker
					//randomly add broker to another broker
					randIndex := 0
					addTo = Brokers[randIndex]
					for i := 0; i < len(Brokers); i++ {
						if len(addTo.Brokers) < MAX_NEIGHBOR {
							// random
							// random int
							randIndex = rand.Intn(len(Brokers))
							addTo = Brokers[randIndex]
							if addTo.Port == broker.Port {
								continue
							}
							break
						}
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
			}

			for _, broker := range Brokers {
				// show neighbor
				fmt.Printf("broker %v has %v neighbor\n", broker.Id, len(broker.Brokers))
				for _, neighbor := range broker.Brokers {
					fmt.Printf("neighbor: %v\n", neighbor.Port)
				}
			}
			// close the goroutine
			return
		}
	}()

	_cli.SeedCliLoop(rpcBrokerClient, Brokers)
}
