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
	RPCAddBroker(ip, port, myId, myIp, myPort string) (*pb.AddBrokerResponse, error)
	RPCAddPublisher(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error)
	RPCAddSubscriber(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error)

	RPCSendAdvertisement(ip, port, myId, myIp, myPort, subject, operator, value string, hopCount int64, nodeType string) (*pb.SendMessageResponse, error)
	RPCSendSubscription(ip, port, myId, myIp, myPort, subject, operator, value, nodeType string) (*pb.SendMessageResponse, error) // hopCount int64
	RPCSendPublication(ip, port, myId, myIp, myPort, subject, operator, value, nodeType string) (*pb.SendMessageResponse, error)

	// RPCSendMessageToBroker(ip string, port int, message string) error
}

func NewBrokerClient() BrokerClient {
	return &brokerClient{}
}

func (bc *brokerClient) RPCAddBroker(ip, port, myId, myIp, myPort string) (*pb.AddBrokerResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddBroker(ctx, &pb.AddBrokerRequest{
		Id:   myId,
		Ip:   myIp,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("Response Error: %v", err)
		return &pb.AddBrokerResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCAddPublisher(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddPublisher(ctx, &pb.AddClientRequest{
		Id:   myId,
		Ip:   myIp,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("Response Error: %v", err)
		return &pb.AddClientResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCAddSubscriber(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.AddSubscriber(ctx, &pb.AddClientRequest{
		Id:   myId,
		Ip:   myIp,
		Port: myPort,
	})
	if err != nil {
		log.Fatalf("Response Error: %v", err)
		return &pb.AddClientResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendAdvertisement(ip, port, myId, myIp, myPort, subject, operator, value string, hopCount int64, nodeType string) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.SendAdvertisement(ctx, &pb.SendMessageRequest{
		Id:       myId,
		Ip:       myIp,
		Port:     myPort,
		Subject:  subject,
		Operator: operator,
		Value:    value,
		HopCount: hopCount,
		NodeType: nodeType,
	})
	if err != nil {
		log.Fatalf("Response Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendSubscription(ip, port, myId, myIp, myPort, subject, operator, value, nodeType string) (*pb.SendMessageResponse, error) { // hopCount int64
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.SendSubscription(ctx, &pb.SendMessageRequest{
		Id:       myId,
		Ip:       myIp,
		Port:     myPort,
		Subject:  subject,
		Operator: operator,
		Value:    value,
		NodeType: nodeType,
		// HopCount: hopCount,
	})
	if err != nil {
		log.Fatalf("Response Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

func (bc *brokerClient) RPCSendPublication(ip, port, myId, myIp, myPort, subject, operator, value, nodeType string) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := rpcConnectTo(ipAddr)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := c.SendPublication(ctx, &pb.SendMessageRequest{
		Id:       myId,
		Ip:       myIp,
		Port:     myPort,
		Subject:  subject,
		Operator: operator,
		Value:    value,
		// HopCount: hopCount,
	})
	if err != nil {
		log.Fatalf("Response Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response, nil
}

// Conncet to gprc server
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
