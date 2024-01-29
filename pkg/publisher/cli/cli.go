package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"wowsan/constants"
	grpcClient "wowsan/pkg/broker/transport"
	model "wowsan/pkg/model"
)

func ExecutionLoop(ip, port string) {
	// defaultIP := "localhost"
	log.Printf("Interactive shell")
	log.Printf("Commands: add, adv")
	id := ip + ":" + port

	publisherModel := model.NewPublisher(id, ip, port)

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
		case "adv":
			if len(args) != 6 {
				//sendAdv apple > 100 localhost 55122
				fmt.Printf("Usage: adv <sbj> <op> <val> <ip> <port>\n")
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
			hopCount := int64(0)

			if publisherModel.Broker == nil {
				response, err := rpcClient.RPCAddPublisher(brokerIp, brokerPort, myId, myIp, myPort)
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				publisherModel.SetBroker(response.Id, response.Ip, response.Port)
				fmt.Printf("Set broker: %s %s %s\n", response.Id, response.Ip, response.Port)
			}

			rpcClient.RPCSendAdvertisement(
				publisherModel.Broker.Ip,
				publisherModel.Broker.Port,
				publisherModel.Id,
				publisherModel.Ip,
				publisherModel.Port,
				subject,
				operator,
				value,
				hopCount,
				constants.PUBLISHER,
			)

		case "pub":
			if len(args) != 6 {
				//sendAdv apple > 100 localhost 55122
				fmt.Printf("Usage: pub <sbj> <op> <val> <ip> <port>\n")
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

			if publisherModel.Broker == nil {
				fmt.Printf("This publisher isn't registered to a broker.\n")
			}
			if publisherModel.Broker.Ip != brokerIp || publisherModel.Broker.Port != brokerPort {
				fmt.Printf("This publisher isn't registered to the broker.\n")
			}

			rpcClient.RPCSendPublication(
				brokerIp,
				brokerPort,
				myId,
				myIp,
				myPort,
				subject,
				operator,
				value,
				constants.PUBLISHER,
			)

		case "broker":
			fmt.Println(publisherModel.Broker.Id, publisherModel.Broker.Ip, publisherModel.Broker.Port)

		}
	}
}
