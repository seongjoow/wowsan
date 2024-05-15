package main

import (
	"flag"
	"fmt"
	cli "wowsan/pkg/broker/cli"
	client "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/broker/service"
	model "wowsan/pkg/model"
)

func main() {
	// ip := flag.String("ip", "", "ID of this node.")
	ip := "localhost"
	port := flag.String("port", "", "Port that this node should listen on.")
	dirIndex := flag.String("dir_index", "", "Directory index for log files.")
	flag.Parse()

	brokerService := service.NewBrokerService(ip, *port, *dirIndex)
	brokerClient := client.NewBrokerClient()

	brokersToAdd := []model.Broker{}

	switch *port {
	case "50001":
		brokerToAdd := &model.Broker{
			Id:   "2",
			Ip:   "localhost",
			Port: "50002",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50002":
		brokerToAdd := &model.Broker{
			Id:   "3",
			Ip:   "localhost",
			Port: "50003",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

		brokerToAdd = &model.Broker{
			Id:   "6",
			Ip:   "localhost",
			Port: "50006",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50003":
		brokerToAdd := &model.Broker{
			Id:   "4",
			Ip:   "localhost",
			Port: "50004",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

		brokerToAdd = &model.Broker{
			Id:   "6",
			Ip:   "localhost",
			Port: "50006",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50004":
		brokerToAdd := &model.Broker{
			Id:   "5",
			Ip:   "localhost",
			Port: "50005",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

		brokerToAdd = &model.Broker{
			Id:   "8",
			Ip:   "localhost",
			Port: "50008",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50006":
		brokerToAdd := &model.Broker{
			Id:   "7",
			Ip:   "localhost",
			Port: "50007",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50008":
		brokerToAdd := &model.Broker{
			Id:   "9",
			Ip:   "localhost",
			Port: "50009",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

		brokerToAdd = &model.Broker{
			Id:   "10",
			Ip:   "localhost",
			Port: "50010",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)
	}

	for _, brokerToAdd := range brokersToAdd {

		broker := brokerService.GetBroker()

		_, err := brokerClient.RPCAddBroker(
			brokerToAdd.Ip,
			brokerToAdd.Port,
			broker.Id,
			broker.Ip,
			broker.Port,
		)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		brokerService.AddBroker(brokerToAdd.Ip, brokerToAdd.Port)
	}

	// brokerToAdd := service.NewBrokerService("localhost", "50002", *dirIndex)
	// brokerToAdd := &model.Broker{
	// 	Id:   "2",
	// 	Ip:   "localhost",
	// 	Port: "50002",
	// }

	// cli.ExecutionLoop(*ip, *port)
	cli.ExecutionLoop(brokerService)
}
