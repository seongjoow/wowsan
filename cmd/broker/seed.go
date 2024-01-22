package main

import (
	"log"
	"net"
	"wowsan/pkg/broker"
	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/broker"

	grpcClient "wowsan/pkg/broker/transport"

	grpc "google.golang.org/grpc"
)

type SeedBroker struct {
	port string
}

var SeedBrokers = []SeedBroker{
	{"50051"},
	{"50052"},
	{"50053"},
}

var Brokers = []*model.Broker{}

func initSeed(port string) (b *model.Broker) {
	lis, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()

	id := "localhost" + ":" + port
	localBrokerModel := model.NewBroker(id, "localhost", port)
	server := broker.NewBrokerRPCServer(localBrokerModel)

	pb.RegisterBrokerServiceServer(s, server)
	go s.Serve(lis)

	//add to Brokers
	return b
}

func main() {
	var Brokers = []*model.Broker{}
	rpcClient := grpcClient.NewBrokerClient()
	for _, seedBroker := range SeedBrokers {
		broker := initSeed(seedBroker.port)
		Brokers = append(Brokers, broker)
	}

	rpcClient.RPCAddBroker(
		"localhost",                         //remote
		SeedBrokers[1].port,                 //remote
		"localhost"+":"+SeedBrokers[0].port, //my
		"localhost",                         //my
		SeedBrokers[0].port,                 //my
	)

	rpcClient.RPCAddBroker(
		"localhost",
		SeedBrokers[2].port,
		"localhost"+":"+SeedBrokers[1].port,
		"localhost",
		SeedBrokers[1].port,
	)

	rpcClient.RPCAddBroker(
		"localhost",
		SeedBrokers[0].port,
		"localhost"+":"+SeedBrokers[2].port,
		"localhost",
		SeedBrokers[2].port,
	)

	// time.Sleep(30 * time.Second)

	// for {
	// 	for index, item := range Brokers[0].SRT {
	// 		fmt.Printf("Brokers 50051 SRT: %s %s %d\n", item.LastHop[index].ID, item.LastHop[index].NodeType, item.HopCount)
	// 	}
	// 	for index, item := range Brokers[1].SRT {
	// 		fmt.Printf("Brokers 50052 SRT: %s %s %d\n", item.LastHop[index].ID, item.LastHop[index].NodeType, item.HopCount)
	// 	}
	// 	for index, item := range Brokers[2].SRT {
	// 		fmt.Printf("Brokers 50053 SRT: %s %s %d\n", item.LastHop[index].ID, item.LastHop[index].NodeType, item.HopCount)
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }
}

// defaultIP := "localhost"
// log.Printf("Interactive shell")
// log.Printf("Commands: add, show, sendAdv")
// id := ip + ":" + port

// // rpc server
// lis, err := net.Listen("tcp", ip+":"+port)
// if err != nil {
// 	log.Fatalf("failed to listen: %v\n", err)
// }

// s := grpc.NewServer()

// localBrokerModel := model.NewBroker(id, ip, port)
// server := broker.NewBrokerRPCServer(localBrokerModel)

// pb.RegisterBrokerServiceServer(s, server)

// go s.Serve(lis)
// // if err := s.Serve(lis); err != nil {
// // 	log.Fatalf("failed to serve: %v\n", err)
// // }
// fmt.Printf("Broker server listening at %v\n", lis.Addr())

// // rpc client
// rpcClient := grpcClient.NewBrokerClient()
// scanner := bufio.NewScanner(os.Stdin)
