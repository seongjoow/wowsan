// broker.go
package broker

import (
	// "fmt"
	"fmt"
	model "wowsan/pkg/model"

	"context"
	grpcClient "wowsan/pkg/broker/transport"
	pb "wowsan/pkg/proto/broker"
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
		Id:   brokerRpcServer.brokerModel.Id,
		Ip:   brokerRpcServer.brokerModel.Ip,
		Port: brokerRpcServer.brokerModel.Port,
	}, nil
}

func (brokerRpcServer *brokerRPCServer) AddPublisher(ctx context.Context, request *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	brokerRpcServer.brokerModel.AddPublisher(request.Id, request.Ip, request.Port)
	fmt.Printf("AddPublisher: %s %s %s\n", request.Id, request.Ip, request.Port)

	return &pb.AddClientResponse{
		Id:   brokerRpcServer.brokerModel.Id,
		Ip:   brokerRpcServer.brokerModel.Ip,
		Port: brokerRpcServer.brokerModel.Port,
	}, nil
}

func (brokerRpcServer *brokerRPCServer) AddSubscriber(ctx context.Context, request *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	brokerRpcServer.brokerModel.AddSubscriber(request.Id, request.Ip, request.Port)
	fmt.Printf("AddSubscriber: %s %s %s\n", request.Id, request.Ip, request.Port)

	return &pb.AddClientResponse{
		Id:   brokerRpcServer.brokerModel.Id,
		Ip:   brokerRpcServer.brokerModel.Ip,
		Port: brokerRpcServer.brokerModel.Port,
	}, nil
}

func (brokerRpcServer *brokerRPCServer) SendAdvertisement(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fmt.Printf("SendAdvertisement: %s %s %s %s %s %s %s %d %s %s\n", request.Id, request.Ip, request.Port, request.Subject, request.Operator, request.Value, request.NodeType, request.HopCount, request.MessageId, request.SenderId)
	if request.Ip == "" {
		return &pb.SendMessageResponse{Message: "IP can't be empty"}, nil
	}

	go brokerRpcServer.brokerModel.PushAdvertisementToQueue(
		&model.AdvertisementRequest{
			Id:          request.Id,
			Ip:          request.Ip,
			Port:        request.Port,
			Subject:     request.Subject,
			Operator:    request.Operator,
			Value:       request.Value,
			NodeType:    request.NodeType,
			HopCount:    request.HopCount,
			MessageId:   request.MessageId,
			PublisherId: request.SenderId,
		},
	)

	// 메세지 큐 구현 전 코드
	// err := brokerRpcServer.brokerModel.SendAdvertisement(
	// 	request.Id,
	// 	request.Ip,
	// 	request.Port,
	// 	request.Subject,
	// 	request.Operator,
	// 	request.Value,
	// 	request.HopCount,
	// 	request.NodeType,
	// )
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// 	return &pb.SendMessageResponse{Message: "fail"}, err
	// }

	return &pb.SendMessageResponse{Message: "success"}, nil
}

func (brokerRpcServer *brokerRPCServer) SendSubscription(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fmt.Printf("SendSubscription: %s %s %s %s %s %s %s\n", request.Id, request.Ip, request.Port, request.Subject, request.Operator, request.Value, request.NodeType)
	if request.Ip == "" {
		return &pb.SendMessageResponse{Message: "IP can't be empty"}, nil
	}

	go brokerRpcServer.brokerModel.PushSubscriptionToQueue(
		&model.SubscriptionRequest{
			Id:       request.Id,
			Ip:       request.Ip,
			Port:     request.Port,
			Subject:  request.Subject,
			Operator: request.Operator,
			Value:    request.Value,
			NodeType: request.NodeType,
		},
	)

	// 메세지 큐 구현 전 코드
	// err := brokerRpcServer.brokerModel.SendSubscription(
	// 	request.Id,
	// 	request.Ip,
	// 	request.Port,
	// 	request.Subject,
	// 	request.Operator,
	// 	request.Value,
	// 	request.NodeType,
	// )
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// 	return &pb.SendMessageResponse{Message: "fail"}, err
	// }

	return &pb.SendMessageResponse{Message: "success"}, nil
}

func (brokerRpcServer *brokerRPCServer) SendPublication(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fmt.Printf("SendPublication: %s %s %s %s %s %s %s\n", request.Id, request.Ip, request.Port, request.Subject, request.Operator, request.Value, request.NodeType)
	if request.Ip == "" {
		return &pb.SendMessageResponse{Message: "IP can't be empty"}, nil
	}

	go brokerRpcServer.brokerModel.PushPublicationToQueue(
		&model.PublicationRequest{
			Id:       request.Id,
			Ip:       request.Ip,
			Port:     request.Port,
			Subject:  request.Subject,
			Operator: request.Operator,
			Value:    request.Value,
			NodeType: request.NodeType,
		},
	)

	// 메세지 큐 구현 전 코드
	// err := brokerRpcServer.brokerModel.SendPublication(
	// 	request.Id,
	// 	request.Ip,
	// 	request.Port,
	// 	request.Subject,
	// 	request.Operator,
	// 	request.Value,
	// 	request.NodeType,
	// )
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// 	return &pb.SendMessageResponse{Message: "fail"}, err
	// }

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
