package transport

// Brocker RPC Client

import (
	"context"
	"fmt"
	"log"
	"time"
	pb "wowsan/pkg/proto"

	grpc "google.golang.org/grpc"
)

type brokerClient struct {
}

type BrokerClient interface {
	RPCAddBroker(ip, port, myId, myIP, myPort string) (*pb.AddBrokerResponse, error)
	RPCSendAdvertisement(ip, port, subject, operator, value, myId, myIP, myPort string, hopCount int64, nodeType string) (*pb.SendMessageResponse, error)
	RPCSendSubscription(ip, port, subject, operator, value, myId, myIP, myPort, nodeType string) (*pb.SendMessageResponse, error) // hopCount int64
	// RPCSendMessageToBroker(ip string, port int, message string) error
}

func NewBrokerClient() BrokerClient {
	return &brokerClient{}
}

func (bc *brokerClient) RPCAddBroker(ip, port, myId, myIP, myPort string) (*pb.AddBrokerResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddBroker(ctx, &pb.AddBrokerRequest{
		Id:   myId,
		Ip:   myIP,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return &pb.AddBrokerResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendAdvertisement(ip, port, subject, operator, value, myId, myIP, myPort string, hopCount int64, nodeType string) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.SendAdvertisement(ctx, &pb.SendMessageRequest{
		Subject:  subject,
		Operator: operator,
		Value:    value,
		Id:       myId,
		Ip:       myIP,
		Port:     myPort,
		HopCount: hopCount,
		NodeType: nodeType,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendSubscription(ip, port, subject, operator, value, myId, myIP, myPort, nodeType string) (*pb.SendMessageResponse, error) { // hopCount int64
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.SendSubscription(ctx, &pb.SendMessageRequest{
		Subject:  subject,
		Operator: operator,
		Value:    value,
		Id:       myId,
		Ip:       myIP,
		Port:     myPort,
		// HopCount: hopCount,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

// func RPCSendMessageToBroker(ip string, port int, message string) error {
// 	ipAddr := ip + ":" + string(port)
// 	ctx, c, cancel, err := connectToBrokerServer(ipAddr)
// 	if err != nil {
// 		log.Fatalf("did not connect: %v\n", err)
// 	}
// 	defer cancel()
// 	_, err = c.SendMessage(ctx, &pb.SendMessageRequest{
// 		Message: message,
// 	})
// 	if err != nil {
// 		log.Fatalf("could not greet: %v", err)
// 	}
// 	return err
// }

// conncet to gprc server
func rpcConnectTo(ip string) (pb.BrokerServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, nil, nil, nil, err
	}

	c := pb.NewBrokerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	return c, conn, ctx, cancel, err
}
