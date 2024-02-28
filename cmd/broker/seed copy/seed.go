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

	return localBrokerModel
}

func main() {
	var isReady bool
	// var isReady = make([]bool, len(seedBrokers))
	var Brokers = []*model.Broker{}
	rpcBrokerClient := grpcClient.NewBrokerClient()

	nodeCount := 10
	seedBrokers := brokerPortsGenerator(nodeCount)
	totalNeighbors := 0
	MAX_NEIGHBOR := 5
	AVG_NEIGHBOR := 2
	MIN_NEIGHBOR := 1

	for index, port := range seedBrokers {
		broker := initSeed(port)
		// isReady[index] = true
		if index == len(seedBrokers)-1 {
			isReady = true
			fmt.Printf("All brokers are ready: %v\n", isReady)
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
					// Randomly add a broker to another broker
					randIndex := 0
					addTo = Brokers[randIndex]

					for i := 0; i < len(Brokers); i++ {
						if len(addTo.Brokers) < MAX_NEIGHBOR {
							randIndex = rand.Intn(len(Brokers))
							addTo = Brokers[randIndex]

							// 자기 자신에게는 추가하지 않음
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
					totalNeighbors += 2
				}
			}

			// 평균 이웃 수에 맞춰 이웃 추가 및 삭제
			// if totalNeighbors < AVG_NEIGHBOR*len(Brokers) {
			for {
				currentAvg := totalNeighbors / len(Brokers)
				fmt.Printf("totalNeighbors / brokers : %v/ %v\n", totalNeighbors, len(Brokers))
				fmt.Printf("currentAvg: %v\n", currentAvg)
				if currentAvg == AVG_NEIGHBOR {
					fmt.Println("currentAvg == AVG_NEIGHBOR (1)")
					break
				}
				if currentAvg < AVG_NEIGHBOR {
					for i := 0; i < len(Brokers); i++ {
						neighborindex := rand.Intn(nodeCount)

						if neighborindex != i && !contains(Brokers[i].Brokers, Brokers[neighborindex]) {
							_, err := rpcBrokerClient.RPCAddBroker(
								Brokers[neighborindex].Ip,
								Brokers[neighborindex].Port,
								Brokers[i].Id,
								Brokers[i].Ip,
								Brokers[i].Port,
							)
							if err != nil {
								fmt.Printf("error: %v", err)
							}
							Brokers[i].AddBroker(Brokers[neighborindex].Id, Brokers[neighborindex].Ip, Brokers[neighborindex].Port)

							totalNeighbors += 2

							currentAvg = totalNeighbors / len(Brokers)
							if currentAvg == AVG_NEIGHBOR {
								fmt.Println("currentAvg == AVG_NEIGHBOR (2)")
								break
							}
						}
					}
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

func contains(brokers map[string]*model.Broker, broker *model.Broker) bool {
	var brokerList []*model.Broker

	for _, broker := range brokers {
		brokerList = append(brokerList, broker)
	}

	for _, b := range brokerList {
		if b.Ip == broker.Ip && b.Port == broker.Port {
			return true
		}
	}
	return false
}
