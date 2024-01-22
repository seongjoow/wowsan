package transport

import (
	"context"
	"log"
	"time"
	pb "wowsan/pkg/proto/subscriber"

	grpc "google.golang.org/grpc"
)

type subscriberClient struct {
}

type SubscriberClient interface {
	RPCReceivePublication(ip, port, subject, operator, value string) (*pb.ReceivePublicationResponse, error)
}

func NewSubscriberClient() SubscriberClient {
	return &subscriberClient{}
}

func (sc *subscriberClient) RPCReceivePublication(ip, port, subject, operator, value string) (*pb.ReceivePublicationResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.ReceivePublication(ctx, &pb.ReceivePublicationRequest{
		Subject:  subject,
		Operator: operator,
		Value:    value,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return &pb.ReceivePublicationResponse{}, err
	}
	return response, nil
}

// rpcConnectTo is a helper function to connect to a remote server
func rpcConnectTo(ipAddr string) (pb.SubscriberServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {
	conn, err := grpc.Dial(ipAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
		return nil, nil, nil, nil, err
	}

	c := pb.NewSubscriberServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	return c, conn, ctx, cancel, nil
}
