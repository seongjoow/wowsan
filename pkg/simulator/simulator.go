package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
	"wowsan/constants"
	grpcClient "wowsan/pkg/broker/transport"
	"wowsan/pkg/model"
	pb "wowsan/pkg/proto/subscriber"
	"wowsan/pkg/subscriber"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// getExpInterval 함수는 지수 분포를 사용하여 다음 호출까지의 대기 시간을 반환함
func getExpInterval(lambda float64) time.Duration {
	expRandom := rand.ExpFloat64() / lambda
	return time.Duration(expRandom * float64(time.Second))
}

// runSimulation 함수는 주어진 시간 동안 시뮬레이션을 실행함
func RunPublisherSimulation(durationSeconds int, lambda float64, brokerIp string, brokerPort string, publisherId string, publisherIp string, publisherPort string) {
	publisherModel := model.NewPublisher(publisherId, publisherIp, publisherPort)
	// rpc client
	rpcClient := grpcClient.NewBrokerClient()

	if publisherModel.Broker == nil {
		response, err := rpcClient.RPCAddPublisher(brokerIp, brokerPort, publisherId, publisherIp, publisherPort)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		publisherModel.SetBroker(response.Id, response.Ip, response.Port)
		fmt.Printf("Set broker: %s %s %s\n", response.Id, response.Ip, response.Port)
	}

	start := time.Now()
	end := start.Add(time.Duration(durationSeconds) * time.Second)

	// 시뮬레이션 루프
	for time.Now().Before(end) {
		subject := "apple"
		operator := ">"
		value := rand.Intn(9999)
		// value to string
		strValue := fmt.Sprintf("%d", value)

		hopCount := int64(0)
		interval := getExpInterval(lambda)
		time.Sleep(interval)

		rpcClient.RPCSendAdvertisement(
			publisherModel.Broker.Ip,
			publisherModel.Broker.Port,
			publisherModel.Id,
			publisherModel.Ip,
			publisherModel.Port,
			subject,
			operator,
			strValue,
			constants.PUBLISHER,
			hopCount,
			uuid.NewV4().String(), // TODO: messageId
			publisherId,
		)
		// 메세지 전달 함수 호출

	}
}

func RunSubscriberSimulation(durationSeconds int, lambda float64, brokerIp string, brokerPort string, subscriberId string, subscriberIp string, subscriberPort string) {
	lis, err := net.Listen("tcp", subscriberIp+":"+subscriberPort)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	subscriberModel := model.NewSubscriber(subscriberId, subscriberIp, subscriberPort)
	// rpc client
	rpcClient := grpcClient.NewBrokerClient()

	server := subscriber.NewSubscriberRPCServer(subscriberModel)
	s := grpc.NewServer()
	pb.RegisterSubscriberServiceServer(s, server)

	go s.Serve(lis)

	if subscriberModel.Broker == nil {
		response, err := rpcClient.RPCAddSubscriber(brokerIp, brokerPort, subscriberId, subscriberIp, subscriberPort)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		subscriberModel.SetBroker(response.Id, response.Ip, response.Port)
		fmt.Printf("Set broker: %s %s %s\n", response.Id, response.Ip, response.Port)
	}

	start := time.Now()
	end := start.Add(time.Duration(durationSeconds) * time.Second)

	// 시뮬레이션 루프
	for time.Now().Before(end) {
		subject := "apple"
		operator := ">"
		value := rand.Intn(9999)
		// value to string
		strValue := fmt.Sprintf("%d", value)

		hopCount := int64(0)
		interval := getExpInterval(lambda)
		time.Sleep(interval)

		rpcClient.RPCSendAdvertisement(
			subscriberModel.Broker.Ip,
			subscriberModel.Broker.Port,
			subscriberModel.Id,
			subscriberModel.Ip,
			subscriberModel.Port,
			subject,
			operator,
			strValue,
			constants.SUBSCRIBER,
			hopCount,
			uuid.NewV4().String(), // TODO: messageId
			subscriberId,
		)
		// 메세지 전달 함수 호출

	}
}

func main() {

	rand.Seed(time.Now().UnixNano())
	lambda := 0.333
	duration := 10

	go RunPublisherSimulation(duration, lambda, "localhost", "50051", "1", "localhost", "1111")
	go RunSubscriberSimulation(duration, lambda, "localhost", "50054", "2", "localhost", "2222")
	// block thread
	select {}
}
