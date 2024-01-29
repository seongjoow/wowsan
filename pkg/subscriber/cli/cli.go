package cli

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"wowsan/constants"
	grpcClient "wowsan/pkg/broker/transport"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/subscriber"
	"wowsan/pkg/subscriber"

	"google.golang.org/grpc"
)

func ExecutionLoop(ip, port string) {
	log.Printf("Interactive shell")
	log.Printf("Commands: add, sub")
	id := ip + ":" + port

	// rpc server
	lis, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	subscriberModel := model.NewSubscriber(id, ip, port)
	server := subscriber.NewSubscriberRPCServer(subscriberModel)

	pb.RegisterSubscriberServiceServer(s, server)

	go s.Serve(lis)
	// rpc client
	rpcClient := grpcClient.NewBrokerClient()
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
		case "sub":
			if len(args) != 6 {
				//sub apple = 80 localhost 55122
				fmt.Printf("Usage: sub <sbj> <op> <val> <ip> <port>\n")
				continue
			}

			subject := args[1]
			operator := args[2]
			value := args[3]
			brokerIp := args[4]
			brokerPort := args[5]
			myId := id
			myIp := ip
			myPort := port

			// if subscriberModel.IsSubscribed(subject, operator, value) {
			// 	fmt.Printf("Already subscribed\n")
			// 	continue
			// }

			if subscriberModel.Broker == nil {
				response, err := rpcClient.RPCAddSubscriber(brokerIp, brokerPort, myId, myIp, myPort)
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				subscriberModel.SetBroker(response.Id, response.Ip, response.Port)
				fmt.Printf("Set broker: %s %s %s\n", response.Id, response.Ip, response.Port)
			}

			rpcClient.RPCSendSubscription(
				subscriberModel.Broker.Ip,
				subscriberModel.Broker.Port,
				subscriberModel.Id,
				subscriberModel.Ip,
				subscriberModel.Port,
				subject,
				operator,
				value,
				constants.SUBSCRIBER,
			)

		case "broker":
			fmt.Println(subscriberModel.Broker.Id, subscriberModel.Broker.Ip, subscriberModel.Broker.Port)

		}
	}
}
