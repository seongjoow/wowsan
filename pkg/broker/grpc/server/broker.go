// broker.go
package server

import (
	"fmt"
	"wowsan/constants"
	model "wowsan/pkg/model"

	"context"
	_brokerUsecase "wowsan/pkg/broker/usecase"
	pb "wowsan/pkg/proto/broker"
)

type brokerRPCServer struct {
	brokerUsecase _brokerUsecase.BrokerUsecase
	pb.UnimplementedBrokerServiceServer
}

func NewBrokerRPCServer(brokerUsecase _brokerUsecase.BrokerUsecase) *brokerRPCServer {
	return &brokerRPCServer{
		brokerUsecase: brokerUsecase,
	}
}

func (brokerRpcServer *brokerRPCServer) AddBroker(ctx context.Context, request *pb.AddBrokerRequest) (*pb.AddBrokerResponse, error) {
	broker, err := brokerRpcServer.brokerUsecase.AddBroker(request.Id, request.Ip, request.Port)
	if err != nil {
		return &pb.AddBrokerResponse{}, err
	}
	fmt.Printf("AddBroker: %s %s %s\n", request.Id, request.Ip, request.Port)
	return &pb.AddBrokerResponse{
		Id:   broker.Id,
		Ip:   broker.Ip,
		Port: broker.Port,
	}, nil
}

func (brokerRpcServer *brokerRPCServer) AddPublisher(ctx context.Context, request *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	_, err := brokerRpcServer.brokerUsecase.AddPublisher(request.Id, request.Ip, request.Port)
	if err != nil {
		return &pb.AddClientResponse{}, err
	}
	fmt.Printf("AddPublisher: %s %s %s\n", request.Id, request.Ip, request.Port)
	broker := brokerRpcServer.brokerUsecase.GetBroker()
	return &pb.AddClientResponse{
		Id:   broker.Id,
		Ip:   broker.Ip,
		Port: broker.Port,
	}, nil
}

func (brokerRpcServer *brokerRPCServer) AddSubscriber(ctx context.Context, request *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	brokerRpcServer.brokerUsecase.AddSubscriber(request.Id, request.Ip, request.Port)
	fmt.Printf("AddSubscriber: %s %s %s\n", request.Id, request.Ip, request.Port)

	broker := brokerRpcServer.brokerUsecase.GetBroker()
	return &pb.AddClientResponse{
		Id:   broker.Id,
		Ip:   broker.Ip,
		Port: broker.Port,
	}, nil
}

func (brokerRpcServer *brokerRPCServer) SendAdvertisement(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fmt.Printf("SendAdvertisement: %s %s %s %s %s %s %s %d %s %s\n", request.Id, request.Ip, request.Port, request.Subject, request.Operator, request.Value, request.NodeType, request.HopCount, request.MessageId, request.SenderId)
	if request.Ip == "" {
		return &pb.SendMessageResponse{Message: "IP can't be empty"}, nil
	}
	performanceInfoArrary := brokerRpcServer.PbPerformanceInfoToModelPerformanceInfo(request.PerformanceInfo)
	go brokerRpcServer.brokerUsecase.PushMessageToQueue(
		&model.MessageRequest{
			Id:              request.Id,
			Ip:              request.Ip,
			Port:            request.Port,
			Subject:         request.Subject,
			Operator:        request.Operator,
			Value:           request.Value,
			NodeType:        request.NodeType,
			HopCount:        request.HopCount,
			MessageId:       request.MessageId,
			SenderId:        request.SenderId,
			MessageType:     constants.ADVERTISEMENT,
			PerformanceInfo: performanceInfoArrary,
		},
	)

	// 메세지 큐 구현 전 코드
	// err := brokerRpcServer.brokerUsecase.SendAdvertisement(
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
	fmt.Printf("SendSubscription: %s %s %s %s %s %s %s %s %s\n", request.Id, request.Ip, request.Port, request.Subject, request.Operator, request.Value, request.NodeType, request.MessageId, request.SenderId)
	if request.Ip == "" {
		return &pb.SendMessageResponse{Message: "IP can't be empty"}, nil
	}

	performanceInfoArrary := brokerRpcServer.PbPerformanceInfoToModelPerformanceInfo(request.PerformanceInfo)
	go brokerRpcServer.brokerUsecase.PushMessageToQueue(
		&model.MessageRequest{
			Id:              request.Id,
			Ip:              request.Ip,
			Port:            request.Port,
			Subject:         request.Subject,
			Operator:        request.Operator,
			Value:           request.Value,
			NodeType:        request.NodeType,
			MessageId:       request.MessageId,
			SenderId:        request.SenderId,
			MessageType:     constants.SUBSCRIPTION,
			PerformanceInfo: performanceInfoArrary,
		},
	)

	// 메세지 큐 구현 전 코드
	// err := brokerRpcServer.brokerUsecase.SendSubscription(
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

	performanceInfoArrary := brokerRpcServer.PbPerformanceInfoToModelPerformanceInfo(request.PerformanceInfo)
	go brokerRpcServer.brokerUsecase.PushMessageToQueue(
		&model.MessageRequest{
			Id:              request.Id,
			Ip:              request.Ip,
			Port:            request.Port,
			Subject:         request.Subject,
			Operator:        request.Operator,
			Value:           request.Value,
			NodeType:        request.NodeType,
			MessageType:     constants.PUBLICATION,
			PerformanceInfo: performanceInfoArrary,
		},
	)

	// 메세지 큐 구현 전 코드
	// err := brokerRpcServer.brokerUsecase.SendPublication(
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

func (brokerRpcServer *brokerRPCServer) PbPerformanceInfoToModelPerformanceInfo(performanceInfo []*pb.PerformanceInfo) []*model.PerformanceInfo {
	performanceInfoArrary := []*model.PerformanceInfo{}
	// pb to model struct
	for _, performanceInfo := range performanceInfo {
		performanceInfoArrary = append(performanceInfoArrary, &model.PerformanceInfo{
			BrokerId:         performanceInfo.BrokerId,
			Cpu:              performanceInfo.Cpu,
			Memory:           performanceInfo.Memory,
			QueueLength:      performanceInfo.QueueLength,
			QueueTime:        performanceInfo.QueueTime,
			ServiceTime:      performanceInfo.ServiceTime,
			ResponseTime:     performanceInfo.ResponseTime,
			InterArrivalTime: performanceInfo.InterArrivalTime,
			Throughput:       performanceInfo.Throughput,
			Timestamp:        performanceInfo.Timestamp,
		})
	}
	return performanceInfoArrary
}
