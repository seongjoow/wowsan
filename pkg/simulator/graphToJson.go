package simulator

import (
	"encoding/json"
	"log"
	"os"
	model "wowsan/pkg/model"
)

// GraphToJSON 함수는 브로커 노드 연결 관계를 JSON 형식으로 반환함
func GraphToJSON(brokers []*model.Broker) {
	// 각 브로커의 이웃 브로커들을 저장할 맵
	brokerNeighbors := make(map[string][]string)

	for _, broker := range brokers {
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
