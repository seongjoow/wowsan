package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
	"wowsan/pkg/broker"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/broker"

	grpcClient "wowsan/pkg/broker/transport"

	"github.com/sirupsen/logrus"
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

func performanceLogger(
	broker *model.Broker,
	logger *logrus.Logger) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		queueLength := len(broker.MessageQueue)
		queueTime := broker.QueueTime
		serviceTime := broker.ServiceTime
		var throughput float64
		if queueTime+serviceTime == 0 {
			throughput = 0
		} else {
			throughput = 1e9 / float64(queueTime+serviceTime) // 초 단위의 값을 얻기 위해서는 나노초 값을 초로 변환 (time.Duration은 기본적으로 나노초 단위의 정수값을 가짐)
		}

		interArrivalTime := broker.InterArrivalTime

		logger.Printf("Queue length: %v, Queue time: %v, Service time: %v, Throughput: %v, Inter-Arrival Time: %v\n", queueLength, queueTime, serviceTime, throughput, interArrivalTime)
	}
}

func initSeed(port string) *model.Broker {
	lis, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()

	id := "localhost" + ":" + port
	l, err := logger.NewLogger(port)
	if err != nil {
		log.Fatalf("failed to create logger: %v\n", err)
	}
	localBrokerModel := model.NewBroker(id, "localhost", port, l)
	server := broker.NewBrokerRPCServer(localBrokerModel)

	pb.RegisterBrokerServiceServer(s, server)
	go s.Serve(lis)
	go localBrokerModel.DoMessageQueue()
	go performanceLogger(localBrokerModel, l)
	l.WithFields(logrus.Fields{
		"port": port,
	}).Info("Broker server listening")

	fmt.Printf("Broker server listening at %v\n", lis.Addr())

	return localBrokerModel
}

func main() {
	var isReady bool
	var Brokers = []*model.Broker{}
	rpcBrokerClient := grpcClient.NewBrokerClient()

	nodeCount := 10
	seedBrokers := brokerPortsGenerator(nodeCount)
	totalNeighbors := 0
	MAX_NEIGHBOR := 5
	AVG_NEIGHBOR := 3
	MIN_NEIGHBOR := 1
	var brokerToAdd *model.Broker

	for index, port := range seedBrokers {
		broker := initSeed(port)
		if index == len(seedBrokers)-1 {
			isReady = true
			fmt.Printf("All brokers are ready: %v\n", isReady)
		}
		Brokers = append(Brokers, broker)
	}

	go func() {
		if isReady {
			for _, broker := range Brokers {
				if len(broker.Brokers) > MAX_NEIGHBOR {
					break
				} else if len(broker.Brokers) < MIN_NEIGHBOR {

					// Randomly add a broker to another broker
					for len(broker.Brokers) < MIN_NEIGHBOR {
						randIndex := rand.Intn(nodeCount)
						brokerToAdd = Brokers[randIndex]
						if len(brokerToAdd.Brokers) >= MAX_NEIGHBOR {
							continue
						}
						// 자기 자신에게는 추가하지 않음
						if brokerToAdd.Ip == broker.Ip && brokerToAdd.Port == broker.Port {
							continue
						}
						rpcBrokerClient.RPCAddBroker(
							brokerToAdd.Ip,
							brokerToAdd.Port,
							broker.Id,
							broker.Ip,
							broker.Port,
						)
						broker.AddBroker(brokerToAdd.Id, brokerToAdd.Ip, brokerToAdd.Port)
						totalNeighbors += 2
					}
				}
			}

			// 평균 이웃 수에 맞춰 이웃 추가
			for {
				currentAvg := totalNeighbors / nodeCount
				fmt.Printf("totalNeighbors / brokers : %v/ %v\n", totalNeighbors, nodeCount)
				fmt.Printf("currentAvg: %v\n", currentAvg)
				if currentAvg == AVG_NEIGHBOR {
					fmt.Println("currentAvg == AVG_NEIGHBOR (1)")
					break
				}
				if currentAvg < AVG_NEIGHBOR {
					for i := 0; i < len(Brokers); i++ {
						randIndex := rand.Intn(nodeCount)
						brokerToAdd = Brokers[randIndex]

						if randIndex != i && !contains(Brokers[i].Brokers, brokerToAdd) {
							_, err := rpcBrokerClient.RPCAddBroker(
								brokerToAdd.Ip,
								brokerToAdd.Port,
								Brokers[i].Id,
								Brokers[i].Ip,
								Brokers[i].Port,
							)
							if err != nil {
								fmt.Printf("error: %v", err)
							}
							Brokers[i].AddBroker(brokerToAdd.Id, brokerToAdd.Ip, brokerToAdd.Port)
							totalNeighbors += 2

							currentAvg = totalNeighbors / nodeCount
							if currentAvg == AVG_NEIGHBOR {
								fmt.Println("currentAvg == AVG_NEIGHBOR (2)")
								break
							}
						}
					}
				}
			}

			for _, broker := range Brokers {
				// Show neighbor

				fmt.Printf("broker %v has %v neighbor\n", broker.Id, len(broker.Brokers))
				for _, neighbor := range broker.Brokers {
					fmt.Printf("neighbor: %v\n", neighbor.Port)
				}
			}
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
