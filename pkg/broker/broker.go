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
	fmt.Printf("SendAdvertisement: %s %s %s %s %s %s %d\n", request.Subject, request.Operator, request.Value, request.Id, request.Ip, request.Port, request.HopCount)
	srt := brokerRPCServer.brokerModel.SRT

	reqSrtItem := model.NewSRTItem(
		request.Subject,
		request.Operator,
		request.Value,
		request.Ip,
		request.Port,
		request.Id,
		request.HopCount,
	)
	hopCount := request.HopCount

	isExist := false

	// TODO
	// Publisher로부터 몇 hop 건너온 메시지인지 확인하는 로직 필요

	// 같은 advertisement에 대한 last hop이 이미 존재하는 경우,
	// 건너온 hop이 더 짧거나 같은 경우에만 last hop 업데이트
	for _, item := range srt {
		if item.Advertisement.Subject == reqSrtItem.Advertisement.Subject &&
			item.Advertisement.Operator == reqSrtItem.Advertisement.Operator &&
			item.Advertisement.Value == reqSrtItem.Advertisement.Value {
			if reqSrtItem.HopCount <= hopCount {
				reqSrtItem.HopCount = hopCount
				if item.HopCount == reqSrtItem.HopCount {
					item.AddLastHop(request.Id, request.Ip, request.Port)
					break
				}
				if reqSrtItem.HopCount < hopCount {
					item = reqSrtItem
					break
				}
			}
			isExist = true
		}
	}

	// 동일한 advertisement가 존재하지 않는 경우 (새로운 advertisement인 경우)
	if isExist == false {
		srt = append(srt, reqSrtItem)
	}

	// brokerRPCServer.brokerModel.SRT = srt // 포인터?로 반환하면 필요 없음?

	// TODO
	// RPC Client
	// rpcClient := grpcClient.NewBrokerClient()
	// advertisement가 온 곳으로부터 멀어지는 방향으로 이웃 브로커들에게 advertisement 전달

	newRequest := &pb.SendMessageRequest{
		Subject:  request.Subject,
		Operator: request.Operator,
		Value:    request.Value,
		Id:       brokerRPCServer.brokerModel.ID,
		Ip:       brokerRPCServer.brokerModel.IP,
		Port:     brokerRPCServer.brokerModel.Port,
		HopCount: request.HopCount + 1,
	}

	for _, neighbor := range brokerRPCServer.brokerModel.Brokers {
		// 온 방향으로는 전송하지 않음
		if neighbor.ID == request.Id {
			continue
		}

		// 새로운 요청을 이웃에게 전송
		_, err := brokerRPCServer.rpcClient.RPCSendAdvertisement(
			neighbor.IP,   //remote broker ip
			neighbor.Port, //remote broker port
			newRequest.Subject,
			newRequest.Operator,
			newRequest.Value,
			newRequest.Id,
			newRequest.Ip,
			newRequest.Port,
			newRequest.HopCount,
		)
		if err != nil {
			log.Fatalf("error: %v", err)
			continue
		}
	}
	return &pb.SendMessageResponse{Message: "success"}, nil

	// item := &model.SubscriptionRoutingTableItem{
	// 	Advertisement: model.Advertisement{
	// 		Key:      request.Subject,
	// 		Value:    request.Value,
	// 		Operator: request.Operator,
	// 	},
	// 	LastHop: []model.LastHop{
	// 		Id:   request.Id,
	// 		Ip:   request.Ip,
	// 		Port: request.Port,
	// 	},
	// }
	// // srt.
	// fmt.Print("SendAdvertisement: ", request.Adv)
	// return &pb.SendAdvertisementResponse{
	// 	Adv: brokerRPCServer.brokerModel.Adv,
	// }, nil
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
