// main.go
package main

import (
	"fmt"
	"log"
	"net"
	"wowsan/pkg/broker"      // 브로커 패키지 경로는 프로젝트에 따라 다를 수 있습니다.
	model "wowsan/pkg/model" // 가정된 모델 패키지

	// "wowsan/pkg/network"     // 네트워크 패키지 경로는 프로젝트에 따라 다를 수 있습니다.
	pb "wowsan/pkg/proto"

	grpc "google.golang.org/grpc"
)

func main() {
	port := ":50051"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	brokerModel := model.NewBroker("1")
	brokerModel.AddBroker("2", "localhost", "50052")

	server := broker.NewBrokerRPCServer(brokerModel)
	pb.RegisterBrokerServiceServer(s, server)
	go s.Serve(lis)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v\n", err)
	// }
	fmt.Printf("Broker server listening at %v\n", lis.Addr())
	select {}
}

// // 네트워크 생성
// net := network.NewNetwork()

// // 브로커 생성 및 네트워크에 추가
// broker := net.NewBroker("broker1")

// // 발행자와 구독자 생성 및 브로커에 등록
// publisher := &broker.Publisher{ID: "pub1"}
// subscriber := &broker.Subscriber{ID: "sub1"}
// broker.RegisterPublisher(publisher)
// broker.RegisterSubscriber(subscriber)

// 추가적인 시스템 초기화 및 실행 로직
// }
