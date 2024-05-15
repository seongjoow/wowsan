package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"wowsan/pkg/broker/service"
)

func ExecutionLoop(service service.BrokerService) {
	log.Printf("Interactive shell")
	log.Printf("Commands: add, srt, prt, broker, publisher, subscriber")

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
				fmt.Println("Usage: add <ip> <port>")
				continue
			}
			service.AddBroker("localhost", args[1])

		case "broker":
			service.Broker()

		case "publisher":
			service.Publisher()

		case "subscriber":
			service.Subscriber()

		case "srt":
			service.Srt()

		case "prt":
			service.Prt()

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
		// 		log.Printf("error: %v", err)
		// 	}
		// 	fmt.Println(res.Message)

		default:
			fmt.Println("Invalid command.")
		}

		fmt.Printf("CMD-> ")
	}
}
