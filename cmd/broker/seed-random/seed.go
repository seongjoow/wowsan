package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	model "wowsan/pkg/model"

	_grpcBrokerClient "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/broker/service"

	_cli "wowsan/cmd/broker/seed-random/cli"
)

// var Brokers = []*model.Broker{}

func brokerPortsGenerator(counts int) (ports []string) {
	for i := 0; i < counts; i++ {
		ports = append(ports, fmt.Sprintf("%d", 50001+i))
	}
	return
}

func main() {
	var isReady bool
	var BrokerServiceList = []service.BrokerService{}
	grpcBrokerClient := _grpcBrokerClient.NewBrokerClient()

	nodeCount := 10
	seedBrokers := brokerPortsGenerator(nodeCount)
	totalNeighbors := 0
	MAX_NEIGHBOR := 4
	AVG_NEIGHBOR := 2
	MIN_NEIGHBOR := 1
	var brokerToAdd *model.Broker

	for index, port := range seedBrokers {
		bService := service.NewBrokerService("localhost", port)
		if index == len(seedBrokers)-1 {
			isReady = true
			fmt.Printf("All brokers are ready: %v\n", isReady)
		}
		BrokerServiceList = append(BrokerServiceList, bService)
	}

	go func() {
		if isReady {
			for _, brokerService := range BrokerServiceList {
				broker := brokerService.GetBroker()

				if len(broker.Brokers) > MAX_NEIGHBOR {

					break
				} else if len(broker.Brokers) < MIN_NEIGHBOR {

					// Randomly add a broker to another broker
					for len(broker.Brokers) < MIN_NEIGHBOR {
						randIndex := rand.Intn(nodeCount)

						brokerToAdd = BrokerServiceList[randIndex].GetBroker()
						if len(brokerToAdd.Brokers) >= MAX_NEIGHBOR {
							continue
						}
						// 자기 자신에게는 추가하지 않음
						if brokerToAdd.Ip == broker.Ip && brokerToAdd.Port == broker.Port {
							continue
						}
						grpcBrokerClient.RPCAddBroker(
							brokerToAdd.Ip,
							brokerToAdd.Port,
							broker.Id,
							broker.Ip,
							broker.Port,
						)
						brokerService.AddBroker(brokerToAdd.Ip, brokerToAdd.Port)
						totalNeighbors += 2
					}
				}
			}

			// 평균 이웃 수에 맞춰 이웃 추가
			for {
				currentAvg := totalNeighbors / nodeCount
				fmt.Printf("totalNeighbors / brokers : %v/ %v\n", totalNeighbors, nodeCount)
				fmt.Printf("currentAvg: %v\n", currentAvg)
				if currentAvg == AVG_NEIGHBOR {
					fmt.Println("currentAvg == AVG_NEIGHBOR (1)")
					break
				}
				if currentAvg < AVG_NEIGHBOR {
					for i := 0; i < len(BrokerServiceList); i++ {
						randIndex := rand.Intn(nodeCount)
						brokerToAdd = BrokerServiceList[randIndex].GetBroker()

						if randIndex != i && !contains(BrokerServiceList[i].GetBroker().Brokers, brokerToAdd) {
							_, err := grpcBrokerClient.RPCAddBroker(
								brokerToAdd.Ip,
								brokerToAdd.Port,
								BrokerServiceList[i].GetBroker().Id,
								BrokerServiceList[i].GetBroker().Ip,
								BrokerServiceList[i].GetBroker().Port,
							)
							if err != nil {
								fmt.Printf("error: %v", err)
							}
							BrokerServiceList[i].AddBroker(brokerToAdd.Ip, brokerToAdd.Port)
							totalNeighbors += 2

							currentAvg = totalNeighbors / nodeCount
							if currentAvg == AVG_NEIGHBOR {
								fmt.Println("currentAvg == AVG_NEIGHBOR (2)")
								break
							}
						}
					}
				}
			}

			// Print neighbors
			// for _, broker := range Brokers {
			//  fmt.Printf("broker %v has %v neighbor\n", broker.Id, len(broker.Brokers))
			//  for _, neighbor := range broker.Brokers {
			//      fmt.Printf("neighbor: %v\n", neighbor.Port)
			//  }
			// }

			// Graph 시각화

			GraphToJSON(BrokerServiceList)

			return
		}
	}()
	_cli.SeedCliLoop(grpcBrokerClient, BrokerServiceList)
}

func contains(brokers map[string]*model.Broker, broker *model.Broker) bool {
	var brokerList []*model.Broker

	for _, broker := range brokers {
		brokerList = append(brokerList, broker)
	}

	for _, b := range brokerList {
		if b.Ip == broker.Ip && b.Port == broker.Port {
			return true
		}
	}
	return false
}

// GraphToJSON 함수는 브로커 노드 연결 관계를 JSON 형식으로 반환함
func GraphToJSON(brokerServiceList []service.BrokerService) {
	// 각 브로커의 이웃 브로커들을 저장할 맵
	brokerNeighbors := make(map[string][]string)

	for _, brokerService := range brokerServiceList {
		broker := brokerService.GetBroker()
		// 브로커의 ID를 키로 하여 이웃 브로커들의 ID를 저장
		for _, neighbor := range broker.Brokers {
			brokerNeighbors[broker.Id] = append(brokerNeighbors[broker.Id], neighbor.Id)
		}
	}

	// JSON 형식으로 변환
	jsonData, err := json.MarshalIndent(brokerNeighbors, "", "    ")
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	// 목표 폴더 경로(wowsan/data/graph) 설정
	folderPath := "data/graph"
	// 필요한 경우 폴더 생성
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %s", err)
	}

	// 파일 저장
	filePath := folderPath + "/graph.json"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error occurred during file creation. Error: %s", err.Error())
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error occurred during writing to file. Error: %s", err.Error())
	}

	log.Println("Graph data has been successfully saved to graph.json")

}
