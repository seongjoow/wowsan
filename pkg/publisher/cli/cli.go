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
	log.Printf("Commands: add, show, sendAdv")
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
		case "sendAdv":
			if len(args) != 6 {
				//sendAdv apple > 100 localhost 55122
				fmt.Printf("Usage: sendAdv <ip> <port>\n")
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
				response, err := rpcClient.RPCAddBroker(brokerIp, brokerPort, myId, myIp, myPort)
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				publisherModel.SetBroker(response.Id, response.Ip, response.Port)
				fmt.Printf("Added broker: %s %s %s\n", response.Id, response.Ip, response.Port)
			}

			rpcClient.RPCSendAdvertisement(
				publisherModel.Broker.IP,
				publisherModel.Broker.Port,
				subject,
				operator,
				value,
				publisherModel.ID,
				publisherModel.IP,
				publisherModel.Port,
				hopCount,
				constants.PUBLISHER,
			)

		}
	}
}
