package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
	cli "wowsan/pkg/broker/cli"
	client "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/broker/service"
	model "wowsan/pkg/model"
)

func main() {
	ip := "localhost"
	port := flag.String("port", "", "Port that this node should listen on.")
	dirIndex := flag.String("dir_index", "", "Directory index for log files.")
	flag.Parse()

	brokerService := service.NewBrokerService(ip, *port, *dirIndex)

	brokerClient := client.NewBrokerClient()
	http.PostForm("http://localhost:8080/init_broker", map[string][]string{"port": {*port}})

	brokersToAdd := []model.Broker{}

	switch *port {
	case "50001":
		brokerToAdd := &model.Broker{
			Id:   "2",
			Ip:   "localhost",
			Port: "50002",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50002":
		brokerToAdd := &model.Broker{
			Id:   "3",
			Ip:   "localhost",
			Port: "50003",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)
		brokerToAdd = &model.Broker{
			Id:   "4",
			Ip:   "localhost",
			Port: "50004",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50003":
		// brokerToAdd := &model.Broker{
		// 	Id:   "4",
		// 	Ip:   "localhost",
		// 	Port: "50004",
		// }
		// brokersToAdd = append(brokersToAdd, *brokerToAdd)

	// 	brokerToAdd = &model.Broker{
	// 		Id:   "6",
	// 		Ip:   "localhost",
	// 		Port: "50006",
	// 	}
	// 	brokersToAdd = append(brokersToAdd, *brokerToAdd)

	case "50004":
		brokerToAdd := &model.Broker{
			Id:   "5",
			Ip:   "localhost",
			Port: "50005",
		}
		brokersToAdd = append(brokersToAdd, *brokerToAdd)
	// case "50006":
	case "50005":
	}

	go func() {
		for {
			resp, _ := http.Get("http://localhost:8080/get_broker_server_ready")
			bodyBytes, _ := io.ReadAll(resp.Body)
			var respJson map[string]interface{} // JSON 응답을 저장할 맵
			err := json.Unmarshal(bodyBytes, &respJson)
			if err != nil {
				// JSON 파싱 에러 처리
				panic(err)
			}
			fmt.Printf("respJson: %v\n", respJson)
			if ready, ok := respJson["ready"].(bool); ok && ready {
				for _, brokerToAdd := range brokersToAdd {

					broker := brokerService.GetBroker()

					_, err := brokerClient.RPCAddBroker(
						brokerToAdd.Ip,
						brokerToAdd.Port,
						broker.Id,
						broker.Ip,
						broker.Port,
					)
					if err != nil {
						fmt.Printf("error: %v\n", err)
					}
					brokerService.AddBroker(brokerToAdd.Ip, brokerToAdd.Port)
				}
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// brokerToAdd := service.NewBrokerService("localhost", "50002", *dirIndex)
	// brokerToAdd := &model.Broker{
	// 	Id:   "2",
	// 	Ip:   "localhost",
	// 	Port: "50002",
	// }

	// cli.ExecutionLoop(*ip, *port)
	cli.ExecutionLoop(brokerService)
}
