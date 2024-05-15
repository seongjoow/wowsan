package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	_grpcBrokerClient "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/broker/service"
	model "wowsan/pkg/model"
)

func findBroker(brokerServiceList []service.BrokerService, port string) *model.Broker {
	for _, brokerService := range brokerServiceList {
		broker := brokerService.GetBroker()
		if broker.Port == port {
			return broker
		}
	}
	fmt.Println("Not found broker")
	return nil
}

func SeedCliLoop(rpcClient _grpcBrokerClient.BrokerClient, brokerServiceList []service.BrokerService) {

	defaultIP := "localhost"
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("CMD-> ")
	for scanner.Scan() {
		fmt.Printf("OUT-> \n")
		line := scanner.Text()
		line = strings.TrimSpace(line)

		args := strings.Split(line, " ")

		if len(args) > 1 {
			args[1] = strings.TrimSpace(args[1])
		}

		if len(args) == 0 {
			continue
		}
		command := args[0]
		switch command {
		case "add":
			if len(args) != 3 {
				fmt.Println("Invalid command.")
				fmt.Println("usage: add [my port] [remote port]")
				continue
			}

			remotePort := args[2]
			remoteIP := defaultIP + ":" + remotePort

			myPort := args[1]
			myId := defaultIP + ":" + myPort
			myIp := defaultIP + ":" + myPort

			response, err := rpcClient.RPCAddBroker(remoteIP, remotePort, myId, myIp, myPort)
			if err != nil {
				log.Printf("error: %v", err)
			}

			for _, brokerService := range brokerServiceList {
				broker := brokerService.GetBroker()
				if broker.Id == myId {
					brokerService.AddBroker(response.Ip, response.Port)
				}
			}

		case "broker":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				fmt.Println("usage: broker [port]")
				continue
			}
			port := args[1]

			broker := findBroker(brokerServiceList, port)
			if broker == nil {
				continue
			}
			for _, broker := range broker.Brokers {
				fmt.Println(broker.Id, broker.Ip, broker.Port)
			}
		case "publisher":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				fmt.Println("usage: publisher [port]")
				continue
			}
			port := args[1]
			broker := findBroker(brokerServiceList, port)
			if broker == nil {
				continue
			}
			for _, publisher := range broker.Publishers {
				fmt.Println(publisher.Id, publisher.Ip, publisher.Port)
			}

		case "subscriber":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				fmt.Println("usage: subscriber [port]")
				continue
			}
			port := args[1]
			broker := findBroker(brokerServiceList, port)
			if broker == nil {
				continue
			}
			for _, subscriber := range broker.Subscribers {
				fmt.Println(subscriber.Id, subscriber.Ip, subscriber.Port)
			}

		case "srt":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				// usage
				fmt.Println("usage: srt [port]")
				continue
			}
			port := args[1]
			broker := findBroker(brokerServiceList, port)
			if broker == nil {
				continue
			}
			for _, item := range broker.SRT {
				fmt.Printf("Adv: %s %s %s\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value)
				for i := 0; i < len(item.LastHop); i++ {
					fmt.Printf("%s | %s | %d\n", item.LastHop[i].Id, item.LastHop[i].NodeType, item.HopCount)
				}
				fmt.Println("----------------------------")
				// fmt.Printf("SRT: %s %s %d\n", item.LastHop[index].ID, item.LastHop[index].NodeType, item.HopCount)
			}

		case "prt":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				continue
			}
			port := args[1]
			broker := findBroker(brokerServiceList, port)

			if broker == nil {
				continue
			}
			for _, item := range broker.PRT {
				fmt.Printf("Sub: %s %s %s\n", item.Subscription.Subject, item.Subscription.Operator, item.Subscription.Value)
				for i := 0; i < len(item.LastHop); i++ {
					fmt.Printf("%s | %s\n", item.LastHop[i].Id, item.LastHop[i].NodeType)
				}
				fmt.Println("----------------------------")
			}

		case "all":
			for _, brokerSerivce := range brokerServiceList {
				broker := brokerSerivce.GetBroker()
				fmt.Println(broker.Id, broker.Ip, broker.Port)
			}
		default:
			fmt.Println("Invalid command.")
		}
		fmt.Printf("CMD-> ")
	}
}
