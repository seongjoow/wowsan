package transport

// Broker RPC Client

import (
	"context"
	"fmt"
	"log"
	"time"
	pb "wowsan/pkg/proto/broker"

	grpc "google.golang.org/grpc"
)

type brokerClient struct {
}

type BrokerClient interface {
	RPCAddBroker(ip, port, myId, myIP, myPort string) (*pb.AddBrokerResponse, error)
	RPCSendAdvertisement(ip, port, subject, operator, value, myId, myIP, myPort string, hopCount int64, nodeType string) (*pb.SendMessageResponse, error)
	RPCSendSubscription(ip, port, subject, operator, value, myId, myIP, myPort, nodeType string) (*pb.SendMessageResponse, error) // hopCount int64
	RPCSendPublication(ip, port, subject, operator, value, myId, myIP, myPort, nodeType string) (*pb.SendMessageResponse, error)

	RPCAddPublisher(ip, port, myId, myIP, myPort string) (*pb.AddClientResponse, error)
	RPCAddSubscriber(ip, port, myId, myIP, myPort string) (*pb.AddClientResponse, error)
	// RPCSendMessageToBroker(ip string, port int, message string) error
}

func NewBrokerClient() BrokerClient {
	return &brokerClient{}
}

func (bc *brokerClient) RPCAddBroker(ip, port, myId, myIP, myPort string) (*pb.AddBrokerResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddBroker(ctx, &pb.AddBrokerRequest{
		Id:   myId,
		Ip:   myIP,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
		return &pb.AddBrokerResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCAddPublisher(ip, port, myId, myIP, myPort string) (*pb.AddClientResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddPublisher(ctx, &pb.AddClientRequest{
		Id:   myId,
		Ip:   myIP,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
		return &pb.AddClientResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCAddSubscriber(ip, port, myId, myIP, myPort string) (*pb.AddClientResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddSubscriber(ctx, &pb.AddClientRequest{
		Id:   myId,
		Ip:   myIP,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
		return &pb.AddClientResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendAdvertisement(ip, port, subject, operator, value, myId, myIP, myPort string, hopCount int64, nodeType string) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
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
		log.Fatalf("Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendSubscription(ip, port, subject, operator, value, myId, myIP, myPort, nodeType string) (*pb.SendMessageResponse, error) { // hopCount int64
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
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
		NodeType: nodeType,
		// HopCount: hopCount,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendPublication(ip, port, subject, operator, value, myId, myIP, myPort, nodeType string) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.SendPublication(ctx, &pb.SendMessageRequest{
		Subject:  subject,
		Operator: operator,
		Value:    value,
		Id:       myId,
		Ip:       myIP,
		Port:     myPort,
		// HopCount: hopCount,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

// conncet to gprc server
func rpcConnectTo(ip string) (pb.BrokerServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Did not connect: %v\n", err)
		return nil, nil, nil, nil, err
	}

	c := pb.NewBrokerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	return c, conn, ctx, cancel, err
}
