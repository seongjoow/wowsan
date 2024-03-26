package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"wowsan/pkg/subscriber/service"
)

func ExecutionLoop(service *service.SubscriberService) {
	log.Printf("Interactive shell")
	log.Printf("Commands: add, sub")

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
				fmt.Printf("Usage: sub <sbj> <op> <val> <ip> <port>\n")
				continue
			}

			subject := args[1]
			operator := args[2]
			value := args[3]
			brokerIp := args[4]
			brokerPort := args[5]
			// myId := id
			// myIp := ip
			// myPort := port

			// if subscriberModel.IsSubscribed(subject, operator, value) {
			// 	fmt.Printf("Already subscribed\n")
			// 	continue
			// }
			service.SubscriberUsecase.Sub(subject, operator, value, brokerIp, brokerPort)
		case "broker":
			fmt.Println(service.Subscriber.Broker.Id, service.Subscriber.Broker.Ip, service.Subscriber.Broker.Port)
		}
	}
}
