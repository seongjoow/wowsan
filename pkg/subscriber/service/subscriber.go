package service

import (
	"log"
	"net"
	_brokerClient "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/model"
	pb "wowsan/pkg/proto/subscriber"
	grpcServer "wowsan/pkg/subscriber/grpc/server"
	_subscriberUsecase "wowsan/pkg/subscriber/usecase"

	"google.golang.org/grpc"
)

type SubscriberService struct {
	Subscriber        *model.Subscriber
	SubscriberUsecase _subscriberUsecase.SubscriberUsecase
	BrokerClient      _brokerClient.BrokerClient
}

func NewSubscriberService(
	ip string,
	port string,
) *SubscriberService {
	lis, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		panic(err)
	}

	id := ip + ":" + port
	subscriber := model.NewSubscriber(id, ip, port)

	brokerClient := _brokerClient.NewBrokerClient()

	subscriberUsecase := _subscriberUsecase.NewSubscriberUsecase(
		subscriber,
		brokerClient,
	)

	s := grpc.NewServer()

	gServer := grpcServer.NewSubscriberRPCServer(subscriberUsecase)
	pb.RegisterSubscriberServiceServer(s, gServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return &SubscriberService{
		Subscriber:        subscriber,
		SubscriberUsecase: subscriberUsecase,
		BrokerClient:      brokerClient,
	}
}
