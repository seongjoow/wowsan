package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"wowsan/pkg/broker"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/broker"

	grpcClient "wowsan/pkg/broker/transport"

	grpc "google.golang.org/grpc"
)

func findBroker(brokers []*model.Broker, port string) *model.Broker {
	for _, broker := range brokers {
		if broker.Port == port {
			return broker
		}
	}
	fmt.Println("Not found broker")
	return nil
}

func SeedCliLoop(rpcClient grpcClient.BrokerClient, brokers []*model.Broker) {

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
				fmt.Println("useage: add [my port] [remote port]")
				continue
			}

			remotePort := args[2]
			remoteIP := defaultIP + ":" + remotePort

			myPort := args[1]
			myId := defaultIP + ":" + myPort
			myIp := defaultIP + ":" + myPort

			response, err := rpcClient.RPCAddBroker(remoteIP, remotePort, myId, myIp, myPort)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			for _, broker := range brokers {
				if broker.Id == myId {
					broker.AddBroker(response.Id, response.Ip, response.Port)
				}
			}

		case "broker":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				fmt.Println("useage: broker [port]")
				continue
			}
			port := args[1]

			broker := findBroker(brokers, port)
			if broker == nil {
				continue
			}
			for _, broker := range broker.Brokers {
				fmt.Println(broker.Id, broker.Ip, broker.Port)
			}
		case "publisher":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				fmt.Println("useage: publisher [port]")
				continue
			}
			port := args[1]
			broker := findBroker(brokers, port)
			if broker == nil {
				continue
			}
			for _, publisher := range broker.Publishers {
				fmt.Println(publisher.Id, publisher.Ip, publisher.Port)
			}

		case "subscriber":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				fmt.Println("useage: subscriber [port]")
				continue
			}
			port := args[1]
			broker := findBroker(brokers, port)
			if broker == nil {
				continue
			}
			for _, subscriber := range broker.Subscribers {
				fmt.Println(subscriber.Id, subscriber.Ip, subscriber.Port)
			}

		case "srt":
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				// useage
				fmt.Println("useage: srt [port]")
				continue
			}
			port := args[1]
			broker := findBroker(brokers, port)
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
			broker := findBroker(brokers, port)

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
			for _, broker := range brokers {
				fmt.Println(broker.Id, broker.Ip, broker.Port)
			}
		default:
			fmt.Println("Invalid command.")
		}
		fmt.Printf("CMD-> ")
	}
}

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
	localBrokerModel := model.NewBroker(id, "localhost", port)
	server := broker.NewBrokerRPCServer(localBrokerModel)

	pb.RegisterBrokerServiceServer(s, server)
	go s.Serve(lis)
	go localBrokerModel.DoAdvertisementQueue()
	go localBrokerModel.DoSubscriptionQueue()
	go localBrokerModel.DoPublicationQueue()

	fmt.Printf("Broker server listening at %v\n", lis.Addr())

	//add to Brokes
	return localBrokerModel
}

func main() {
	//new a isReady array(len = seedBrokers)
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

	SeedCliLoop(rpcBrokerClient, Brokers)
}
