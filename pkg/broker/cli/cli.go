package cli

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"wowsan/pkg/broker"
	grpcClient "wowsan/pkg/broker/transport"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto"

	grpc "google.golang.org/grpc"
)

func ExecutionLoop(ip, port string) {
	defaultIP := "localhost"
	log.Printf("Interactive shell")
	log.Printf("Commands: add, show, sendAdv")
	id := ip + ":" + port

	// rpc server
	lis, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()

	localBrokerModel := model.NewBroker(id, ip, port)
	server := broker.NewBrokerRPCServer(localBrokerModel)

	pb.RegisterBrokerServiceServer(s, server)

	go s.Serve(lis)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v\n", err)
	// }
	fmt.Printf("Broker server listening at %v\n", lis.Addr())

	// rpc client
	rpcClient := grpcClient.NewBrokerClient()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("CMD-> ")
	for scanner.Scan() {
		fmt.Printf("OUT-> \n")
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// args := strings.SplitN(line, " ", 0)
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
			// if len(args) != 3 {
			// 	fmt.Println("Invalid command.")
			// 	continue
			// }
			// remoteIP := args[1]
			// remotePort := args[2]
			if len(args) != 2 {
				fmt.Println("Invalid command.")
				continue
			}
			remoteIP := defaultIP
			remotePort := args[1]
			myId := id
			myIp := ip
			myPort := port
			response, err := rpcClient.RPCAddBroker(remoteIP, remotePort, myId, myIp, myPort)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			localBrokerModel.AddBroker(response.Id, response.Ip, response.Port)
		case "show":
			for _, broker := range localBrokerModel.Brokers {
				fmt.Println(broker.ID, broker.IP, broker.Port)
			}

		case "srt":
			fmt.Println("[SRT]")
			fmt.Println("------------SRT------------")
			for _, item := range localBrokerModel.SRT {
				fmt.Printf("Adv: %s %s %s\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value)
				for i := 0; i < len(item.LastHop); i++ {
					fmt.Printf("%s | %s | %d\n", item.LastHop[i].ID, item.LastHop[i].NodeType, item.HopCount)
				}
				fmt.Println("----------------------------")
				// fmt.Printf("SRT: %s %s %d\n", item.LastHop[index].ID, item.LastHop[index].NodeType, item.HopCount)
			}

		case "prt":
			for index, item := range localBrokerModel.PRT {
				fmt.Printf("PRT: %s %s \n", item.LastHop[index].ID, item.LastHop[index].NodeType)
			}
		// case "sendAdv":
		// 	if len(args) != 6 {
		// 		fmt.Println("Invalid command.")
		// 		continue
		// 	}
		// 	subject := args[1]z
		// 	operator := args[2]
		// 	value := args[3]
		// 	brokerIp := args[4]
		// 	brokerPort := args[5]
		// 	myId := id
		// 	myIp := ip
		// 	myPort := port
		// 	hopCount := int64(0)

		// 	fmt.Println("sendAdv", subject, operator, value, brokerIp, brokerPort, myId, myIp, myPort, hopCount)

		// 	res, err := rpcClient.RPCSendAdvertisement(
		// 		brokerIp,
		// 		brokerPort,
		// 		subject,
		// 		operator,
		// 		value,
		// 		myId,
		// 		myIp,
		// 		myPort,
		// 		hopCount,
		// 	)

		// 	if err != nil {
		// 		log.Fatalf("error: %v", err)
		// 	}
		// 	fmt.Println(res.Message)

		default:
			fmt.Println("Invalid command.")
		}
		fmt.Printf("CMD-> ")
	}
}
