// broker.go
package broker

import (
	// "fmt"
	// "wowsan/pkg/network" // 가정된 네트워크 패키지
	"fmt"
	"log"
	model "wowsan/pkg/model" // 가정된 모델 패키지

	"context"
	grpcClient "wowsan/pkg/broker/transport"
	pb "wowsan/pkg/proto"
)

// server is used to implement helloworld.GreeterServer.
type brokerRPCServer struct {
	rpcClient   grpcClient.BrokerClient
	brokerModel *model.Broker
	pb.UnimplementedBrokerServiceServer
}

func NewBrokerRPCServer(brokerModel *model.Broker) *brokerRPCServer {
	rpcClient := grpcClient.NewBrokerClient()

	return &brokerRPCServer{
		rpcClient:   rpcClient,
		brokerModel: brokerModel,
	}
}

func (brokerRpcServer *brokerRPCServer) AddBroker(ctx context.Context, request *pb.AddBrokerRequest) (*pb.AddBrokerResponse, error) {
	brokerRpcServer.brokerModel.AddBroker(request.Id, request.Ip, request.Port)
	fmt.Printf("AddBroker: %s %s %s\n", request.Id, request.Ip, request.Port)
	return &pb.AddBrokerResponse{
		Id:   brokerRpcServer.brokerModel.ID,
		Ip:   brokerRpcServer.brokerModel.IP,
		Port: brokerRpcServer.brokerModel.Port,
	}, nil
}

func (brokerRPCServer *brokerRPCServer) SendAdvertisement(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fmt.Printf("SendAdvertisement: %s %s %s %s %s %s %d %s\n", request.Subject, request.Operator, request.Value, request.Id, request.Ip, request.Port, request.HopCount, request.NodeType)
	if request.Ip == "" {
		return &pb.SendMessageResponse{Message: "IP can't be empty"}, nil
	}
	err := brokerRPCServer.brokerModel.SendAdvertisement(
		request.Id,
		request.Ip,
		request.Port,
		request.Subject,
		request.Operator,
		request.Value,
		request.HopCount,
		request.NodeType,
	)
	if err != nil {
		log.Fatalf("error: %v", err)
		return &pb.SendMessageResponse{Message: "fail"}, err
	}
	return &pb.SendMessageResponse{Message: "success"}, nil
}

// func (s *broker) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {

// 	return &pb.SendMessageResponse{
// 		Message: in.Message,
// 		}, nil
// }

// // interface
// type Broker interface {
// 	RegisterPublisher(publisher *model.Publisher)
// 	RegisterSubscriber(subscriber *model.Subscriber)
// 	ProcessMessages()
// }

// // ProcessMessages 함수는 발행자로부터 메시지를 수신하고 구독자에게 전달하는 로직을 포함합니다.
// func (b *broker) ProcessMessages() {

// 	for {
// 		msg, err := b.Node.ReceiveMessage()
// 		if err != nil {
// 			fmt.Println("Error receiving message:", err)
// 			continue
// 		}

// 		for _, sub := range b.Subscribers {
// 			err := sub.Node.SendMessage(msg)
// 			if err != nil {
// 				fmt.Printf("Error sending message to subscriber %s: %v\n", sub.ID, err)
// 			}
// 		}
// 	}
// }

// // 여기에 추가적인 브로커 관련 로직을 구현...
