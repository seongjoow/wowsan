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
	log.Printf("Interactive shell")
	log.Printf("Commands: join, rings, show, resource")
	id := ip + ":" + port

	//

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
	// TODO

	fmt.Printf("CMD-> ")
	for scanner.Scan() {
		fmt.Printf("OUT-> \n")
		line := scanner.Text()
		line = strings.TrimSpace(line)

		args := strings.SplitN(line, " ", 4)
		if len(args) > 1 {
			args[1] = strings.TrimSpace(args[1])
		}

		if len(args) == 0 {
			continue
		}
		command := args[0]
		switch command {
		// TODO
		case "add":
			if len(args) != 3 {
				fmt.Println("Invalid command.")
				continue
			}
			remoteIP := args[1]
			remotePort := args[2]
			myId := id
			myIP := ip
			response, err := rpcClient.RPCAddBroker(remoteIP, remotePort, myId, myIP, port)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			localBrokerModel.AddBroker(response.Id, response.Ip, response.Port)
		case "show":
			for _, broker := range localBrokerModel.Brokers {
				fmt.Println(broker.ID, broker.IP, broker.Port)
			}
		case "srt":

		case "prt":
		case "send":
		default:
			fmt.Println("Invalid command.")
		}
		fmt.Printf("CMD-> ")
	}
}
