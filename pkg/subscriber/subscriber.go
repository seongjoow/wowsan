package subscriber

import (
	"context"
	"fmt"
	"log"

	model "wowsan/pkg/model"
	pb "wowsan/pkg/proto/subscriber"
)

type subscriberRPCServer struct {
	subscriberModel *model.Subscriber
	pb.UnimplementedSubscriberServiceServer
}

func NewSubscriberRPCServer(subscriberModel *model.Subscriber) *subscriberRPCServer {
	return &subscriberRPCServer{
		subscriberModel: subscriberModel,
	}
}

func (subscriberRpcServer *subscriberRPCServer) ReceivePublication(ctx context.Context, request *pb.ReceivePublicationRequest) (*pb.ReceivePublicationResponse, error) {
	fmt.Printf("ReceivePublication: %s %s %s %s\n", request.Id, request.Subject, request.Operator, request.Value)
	if request.Id == "" {
		return &pb.ReceivePublicationResponse{Message: "Id can't be empty"}, nil
	}

	err := subscriberRpcServer.subscriberModel.ReceivePublication(
		request.Id,
		request.Subject,
		request.Operator,
		request.Value,
	)
	if err != nil {
		log.Fatalf("error: %v", err)
		return &pb.ReceivePublicationResponse{Message: "fail"}, err
	}

	return &pb.ReceivePublicationResponse{Message: "success"}, nil
}
