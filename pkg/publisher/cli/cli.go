package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"wowsan/pkg/publisher/service"
)

func ExecutionLoop(service *service.PublisherService) {
	log.Printf("Interactive shell")
	log.Printf("Commands: adv, pub, broker")

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
				fmt.Printf("Usage: adv <sbj> <op> <val> <ip> <port>\n")
				continue
			}
			service.PublisherUsecase.Adv(args[1], args[2], args[3], args[4], args[5])

		case "pub":
			if len(args) != 6 {
				fmt.Printf("Usage: pub <sbj> <op> <val> <ip> <port>\n")
				continue
			}
			service.PublisherUsecase.Pub(args[1], args[2], args[3], args[4], args[5])

		case "broker":
			broker := service.Publisher.Broker
			fmt.Printf("Broker: %s %s %s\n", broker.Id, broker.Ip, broker.Port)

		default:
			fmt.Println("Invalid command.")
		}

		fmt.Printf("CMD-> ")
	}
}
