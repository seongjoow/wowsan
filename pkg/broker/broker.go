// broker.go
package broker

import (
	// "fmt"
	// "wowsan/pkg/network" // 가정된 네트워크 패키지
	"fmt"
	model "wowsan/pkg/model" // 가정된 모델 패키지

	"context"
	pb "wowsan/pkg/proto"
)

// server is used to implement helloworld.GreeterServer.
type brokerRPCServer struct {
	brokerModel *model.Broker
	pb.UnimplementedBrokerServiceServer
}

func NewBrokerRPCServer(brokerModel *model.Broker) *brokerRPCServer {
	return &brokerRPCServer{
		brokerModel: brokerModel,
	}
}

func (brokerRpcServer *brokerRPCServer) AddBroker(ctx context.Context, request *pb.AddBrokerRequest) (*pb.AddBrokerResponse, error) {
	brokerRpcServer.brokerModel.AddBroker(request.Id, request.Ip, request.Port)
	fmt.Print("AddBroker: ", request.Id, request.Ip, request.Port)
	return &pb.AddBrokerResponse{
		Id:   brokerRpcServer.brokerModel.ID,
		Ip:   brokerRpcServer.brokerModel.IP,
		Port: brokerRpcServer.brokerModel.Port,
	}, nil
}

func (brokerRPCServer *brokerRPCServer) SendAdvertisement(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	srt := brokerRPCServer.brokerModel.SRT
	srtItem := model.NewSRTItem()
	srtItem.SetAdvertisement(request.Subject, request.Value, request.Operator)

	// TODO
	// Publisher로부터 몇 hop 건너온 메시지인지 확인하는 로직 필요
	// 같은 advertisement에 대한 last hop이 이미 존재하는 경우,
	// 건너온 hop이 더 짧거나 같은 경우에만 last hop에 추가

	srtItem.AddLastHop(request.Id, request.Ip, request.Port)
	srt = append(srt, srtItem)

	// brokerRPCServer.brokerModel.SRT = srt

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
