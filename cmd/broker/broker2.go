package main

import (
	grpcClient "wowsan/pkg/broker/transport"

	model "wowsan/pkg/model"

	"fmt"
	"log"
	"net" // 브로커 패키지 경로는 프로젝트에 따라 다를 수 있습니다.

	// "wowsan/pkg/network"     // 네트워크 패키지 경로는 프로젝트에 따라 다를 수 있습니다.
	"wowsan/pkg/broker"
	pb "wowsan/pkg/proto"

	grpc "google.golang.org/grpc"
)

func main() {
	//flag id, ip, port

	// id := flag.String("id", "1", "id")
	// ip := flag.String("ip", "localhost", "ip")
	// port := flag.String("port", "50052", "port")

	id := "2"
	port := ":50052"
	remoteIP := "localhost"
	remotePort := "50051"

	// rpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()

	localBrokerModel := model.NewBroker(id, "localhost", port)
	server := broker.NewBrokerRPCServer(localBrokerModel)

	pb.RegisterBrokerServiceServer(s, server)

	go s.Serve(lis)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v\n", err)
	// }
	fmt.Printf("Broker server listening at %v\n", lis.Addr())

	// rpc client
	rpcClient := grpcClient.NewBrokerClient()
	response, err := rpcClient.RPCAddBroker(remoteIP, remotePort, localBrokerModel.ID, localBrokerModel.IP, brokerModel.Port)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	localBrokerModel.AddBroker(response.Id, response.Ip, response.Port)

	// _Cli := cli.NewCLI(broker, grpcClient.NewBrokerClient())
	// _Cli.Run()

	// rpcClient := grpcClient.NewBrokerClient()
	// rpcClient.RPCAddBroker("localhost", "50051", "2", "localhost", "50052")

}
