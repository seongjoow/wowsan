package server

import (
	"fmt"

	"context"
	pb "wowsan/pkg/proto/subscriber"
	_subscriberUsecase "wowsan/pkg/subscriber/usecase"
)

type subscriberRPCServer struct {
	subscriberUsecase _subscriberUsecase.SubscriberUsecase
	pb.UnimplementedSubscriberServiceServer
}

func NewSubscriberRPCServer(subscriberUsecase _subscriberUsecase.SubscriberUsecase) *subscriberRPCServer {
	return &subscriberRPCServer{
		subscriberUsecase: subscriberUsecase,
	}
}

func (subscriberRpcServer *subscriberRPCServer) ReceivePublication(ctx context.Context, request *pb.ReceivePublicationRequest) (*pb.ReceivePublicationResponse, error) {
	fmt.Printf("ReceivePublication: %s %s %s %s\n", request.Id, request.Subject, request.Operator, request.Value)
	if request.Id == "" {
		return &pb.ReceivePublicationResponse{Message: "Id can't be empty"}, nil
	}

	err := subscriberRpcServer.subscriberUsecase.ReceivePublication(
		request.Id,
		request.Subject,
		request.Operator,
		request.Value,
	)
	if err != nil {
		fmt.Printf("error: %v", err)
		return &pb.ReceivePublicationResponse{Message: "fail"}, err
	}

	return &pb.ReceivePublicationResponse{Message: "success"}, nil
}
