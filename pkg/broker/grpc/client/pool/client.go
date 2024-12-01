package transport_pool

// Broker RPC Client

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
	"wowsan/pkg/model"
	pb "wowsan/pkg/proto/broker"

	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

type brokerClient struct {
	log      *logrus.Logger
	connPool *ConnectionPool
}

type BrokerClient interface {
	RPCAddBroker(ip, port, myId, myIp, myPort string) (*pb.AddBrokerResponse, error)
	RPCAddPublisher(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error)
	RPCAddSubscriber(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error)

	RPCSendAdvertisement(ip, port, myId, myIp, myPort, subject, operator, value, nodeType string, hopCount int64, messageId, senderId string, performanceMessageArrary []*model.PerformanceInfo) (*pb.SendMessageResponse, error)
	RPCSendSubscription(ip, port, myId, myIp, myPort, subject, operator, value, nodeType, messageId, senderId string, performanceMessageArrary []*model.PerformanceInfo) (*pb.SendMessageResponse, error) // hopCount int64
	RPCSendPublication(ip, port, myId, myIp, myPort, subject, operator, value, nodeType, messageId string, performanceMessageArrary []*model.PerformanceInfo) (*pb.SendMessageResponse, error)

	// RPCSendMessageToBroker(ip string, port int, message string) error
}

func NewBrokerClientWithLogger(log *logrus.Logger) BrokerClient {
	return &brokerClient{
		log: log,
	}
}

func NewBrokerClient() BrokerClient {
	log := logrus.New()
	return &brokerClient{
		log:      log,
		connPool: NewConnectionPool(),
	}
}

func (bc *brokerClient) RPCAddBroker(ip, port, myId, myIp, myPort string) (*pb.AddBrokerResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := bc.rpcConnectTo(ipAddr)
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()
	response, err := bc.RetryCallback(func() (interface{}, error) {
		return c.AddBroker(ctx, &pb.AddBrokerRequest{
			Id:   myId,
			Ip:   myIp,
			Port: myPort,
		})
	}, myPort, port, 3)
	if err != nil {
		log.Printf("Response Error: %v", err)
		return &pb.AddBrokerResponse{}, err
	}
	return response.(*pb.AddBrokerResponse), nil
}

func (bc *brokerClient) RPCAddPublisher(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := bc.rpcConnectTo(ipAddr)
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()
	response, err := bc.RetryCallback(func() (interface{}, error) {
		return c.AddPublisher(ctx, &pb.AddClientRequest{
			Id:   myId,
			Ip:   myIp,
			Port: myPort,
		})
	}, myPort, port, 3)
	if err != nil {
		log.Printf("Response Error: %v", err)
		return &pb.AddClientResponse{}, err
	}
	return response.(*pb.AddClientResponse), nil
}

func (bc *brokerClient) RPCAddSubscriber(ip, port, myId, myIp, myPort string) (*pb.AddClientResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := bc.rpcConnectTo(ipAddr)
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	response, err := bc.RetryCallback(func() (interface{}, error) {
		return c.AddSubscriber(ctx, &pb.AddClientRequest{
			Id:   myId,
			Ip:   myIp,
			Port: myPort,
		})
	}, myPort, port, 3)

	if err != nil {
		log.Printf("Response Error: %v", err)
		return &pb.AddClientResponse{}, err
	}
	return response.(*pb.AddClientResponse), nil
}

func (bc *brokerClient) RPCSendAdvertisement(ip, port, myId, myIp, myPort, subject, operator, value, nodeType string, hopCount int64, messageId, senderId string, performanceInfoArrary []*model.PerformanceInfo) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := bc.rpcConnectTo(ipAddr)
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	pbPerformanceInfoArray := bc.ModelPerformanceInfoToPbPerformanceInfo(performanceInfoArrary)

	response, err := bc.RetryCallback(func() (interface{}, error) {
		return c.SendAdvertisement(ctx, &pb.SendMessageRequest{
			Id:              myId,
			Ip:              myIp,
			Port:            myPort,
			Subject:         subject,
			Operator:        operator,
			Value:           value,
			NodeType:        nodeType,
			HopCount:        hopCount,
			MessageId:       messageId,
			SenderId:        senderId,
			PerformanceInfo: pbPerformanceInfoArray,
		})
	}, myPort, port, 3)
	if err != nil {
		log.Printf("Response Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response.(*pb.SendMessageResponse), nil
}

func (bc *brokerClient) RPCSendSubscription(ip, port, myId, myIp, myPort, subject, operator, value, nodeType, messageId, senderId string, performanceInfoArrary []*model.PerformanceInfo) (*pb.SendMessageResponse, error) { // hopCount int64
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := bc.rpcConnectTo(ipAddr)
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
	}
	defer conn.Close()
	defer cancel()

	pbPerformanceInfoArray := bc.ModelPerformanceInfoToPbPerformanceInfo(performanceInfoArrary)

	response, err := bc.RetryCallback(func() (interface{}, error) {
		return c.SendSubscription(ctx, &pb.SendMessageRequest{
			Id:       myId,
			Ip:       myIp,
			Port:     myPort,
			Subject:  subject,
			Operator: operator,
			Value:    value,
			NodeType: nodeType,
			// HopCount: hopCount,
			MessageId:       messageId,
			SenderId:        senderId,
			PerformanceInfo: pbPerformanceInfoArray,
		})
	}, myPort, port, 3)
	if err != nil {
		log.Printf("Response Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}

	return response.(*pb.SendMessageResponse), nil
}

func (bc *brokerClient) RPCSendPublication(ip, port, myId, myIp, myPort, subject, operator, value, nodeType, messageId string, performanceInfoArray []*model.PerformanceInfo) (*pb.SendMessageResponse, error) {
	ipAddr := ip + ":" + string(port)
	c, conn, ctx, cancel, err := bc.rpcConnectTo(ipAddr)
	if err != nil {
		log.Printf("Did not connect: %v\n", err)
		return &pb.SendMessageResponse{}, err
	}
	defer conn.Close()
	defer cancel()

	pbPerformanceInfoArray := bc.ModelPerformanceInfoToPbPerformanceInfo(performanceInfoArray)

	response, err := bc.RetryCallback(func() (interface{}, error) {
		return c.SendPublication(ctx, &pb.SendMessageRequest{
			Id:              myId,
			Ip:              myIp,
			Port:            myPort,
			Subject:         subject,
			Operator:        operator,
			Value:           value,
			MessageId:       messageId,
			PerformanceInfo: pbPerformanceInfoArray,
		})
	}, myPort, port, 3)
	if err != nil {
		log.Printf("Response Error: %v", err)
		return &pb.SendMessageResponse{}, err
	}
	return response.(*pb.SendMessageResponse), nil
}

// Conncet to gprc server
func (bc *brokerClient) rpcConnectTo(ip string) (pb.BrokerServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {

	// Check if the connection exists in the pool
	conn, err := bc.connPool.GetConnection(ip)
	if err != nil {
		// bc.log.Infof("Connection for %s not found, attempting to create a new one with retry.", ip)

		// if connection doesn't exist, add it to the pool
		_, err = bc.RetryCallback(func() (interface{}, error) {
			err := bc.connPool.AddConnection(ip)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}, "", ip, 3)

		if err != nil {
			fmt.Printf("Failed to add connection after retries: %v\n", err)
			return nil, nil, nil, nil, err
		}

		// Get the connection from the pool
		conn, err = bc.connPool.GetConnection(ip)
		if err != nil {
			fmt.Printf("Failed to get connection after adding: %v\n", err)
			return nil, nil, nil, nil, err
		}
	}

	// Create a new client
	client := pb.NewBrokerServiceClient(conn)
	// bc.clientCache[ip] = client

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	return client, conn, ctx, cancel, nil
}

func (bc *brokerClient) ModelPerformanceInfoToPbPerformanceInfo(performanceInfoArrary []*model.PerformanceInfo) []*pb.PerformanceInfo {
	pbPerformanceInfoArray := []*pb.PerformanceInfo{}
	// model to pb struct
	for _, performanceInfo := range performanceInfoArrary {
		pbPerformanceInfo := &pb.PerformanceInfo{
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
		}
		pbPerformanceInfoArray = append(pbPerformanceInfoArray, pbPerformanceInfo)
	}
	return pbPerformanceInfoArray
}

// Retry Callback function
func (bc *brokerClient) RetryCallback(callback func() (interface{}, error), clientPort, serverPort string, retry int) (interface{}, error) {
	var err error
	var result interface{}
	funcName := runtime.FuncForPC(reflect.ValueOf(callback).Pointer()).Name()

	for i := 0; i < retry; i++ {
		result, err = callback()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		bc.log.WithFields(logrus.Fields{
			"clientPort":   clientPort,
			"serverPort":   serverPort,
			"functionName": funcName,
			"error":        err,
		}).Error("RetryCallback Error")
		return nil, err
	}

	return result, nil
}
