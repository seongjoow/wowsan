package service

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"time"

	_brokerClient "wowsan/pkg/broker/grpc/client"
	grpcServer "wowsan/pkg/broker/grpc/server"
	_brokerUsecase "wowsan/pkg/broker/usecase"
	"wowsan/pkg/logger"
	"wowsan/pkg/model"
	_subscriberClient "wowsan/pkg/subscriber/grpc/client"

	pb "wowsan/pkg/proto/broker"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type BrokerService interface {
	AddBroker(remoteIp, remotePort string)
	Srt()
	Prt()
	Broker()
	Publisher()
	Subscriber()
	GetBroker() *model.Broker
}

type brokerService struct {
	hopLogger        *logrus.Logger
	tickLogger       *logrus.Logger
	brokerUsercase   _brokerUsecase.BrokerUsecase
	brokerClient     _brokerClient.BrokerClient
	subscriberClient _subscriberClient.SubscriberClient
}

func NewBrokerService(
	ip string,
	port string,
	baseDir string,
) BrokerService {
	lis, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		panic(err)
	}

	hopLogger, err := logger.NewLogger(port, "./log/hopLogger/"+baseDir)
	if err != nil {
		panic(err)
	}
	tickLogger, err := logger.NewLogger(port+"_tick", "./log/tickLogger/"+baseDir)
	if err != nil {
		panic(err)
	}

	brokerLoggerDir := "./log/brokerInfoLogger/" + baseDir
	brokerLogger := logger.NewBrokerInfoLogger(brokerLoggerDir)

	id := ip + ":" + port
	broker := model.NewBroker(id, ip, port)

	brokerClient := _brokerClient.NewBrokerClient()
	subscriberClient := _subscriberClient.NewSubscriberClient()

	brokerUsecase := _brokerUsecase.NewBrokerUsecase(
		hopLogger,
		tickLogger,
		brokerLogger,
		broker,
		brokerClient,
		subscriberClient,
	)

	s := grpc.NewServer()

	gServer := grpcServer.NewBrokerRPCServer(brokerUsecase)
	pb.RegisterBrokerServiceServer(s, gServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	go func() {
		for {
			b := brokerUsecase.GetBroker()
			if time.Since(b.Close) > 30*time.Second {
				// terminate the program
				fmt.Println("Broker is closed")
				http.PostForm("http://localhost:8080/done", map[string][]string{"port": {port}})
				os.Exit(1)
			}
		}
	}()
	// api call server init finished

	go brokerUsecase.PerformanceTickLogger(1 * time.Second)
	go brokerUsecase.DoMessageQueue()

	return &brokerService{
		hopLogger:        hopLogger,
		tickLogger:       tickLogger,
		brokerUsercase:   brokerUsecase,
		brokerClient:     brokerClient,
		subscriberClient: subscriberClient,
	}
}

func (b *brokerService) AddBroker(remoteIp, remotePort string) {
	broker := b.brokerUsercase.GetBroker()
	response, err := b.brokerClient.RPCAddBroker(remoteIp, remotePort, broker.Id, broker.Ip, broker.Port)
	if err != nil {
		b.hopLogger.Fatalf("error: %v", err)
	}
	b.brokerUsercase.AddBroker(response.Id, response.Ip, response.Port)
}

func (b *brokerService) Srt() {
	broker := b.brokerUsercase.GetBroker()
	fmt.Println("-----------[SRT]-----------")
	for _, item := range broker.SRT {
		fmt.Printf("Adv: %s %s %s (%s) | %s\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.Identifier.MessageId, item.Identifier.SenderId)
		for i := 0; i < len(item.LastHop); i++ {
			fmt.Printf("%s | %s | %d \n", item.LastHop[i].Id, item.LastHop[i].NodeType, item.HopCount)
		}
		fmt.Println("----------------------------")
		// fmt.Printf("SRT: %s %s %d\n", item.LastHop[index].ID, item.LastHop[index].NodeType, item.HopCount)
	}
}

func (b *brokerService) Prt() {
	broker := b.brokerUsercase.GetBroker()
	fmt.Println("-----------[PRT]-----------")
	for _, item := range broker.PRT {
		fmt.Printf("Sub: %s %s %s (%s) | %s\n", item.Subscription.Subject, item.Subscription.Operator, item.Subscription.Value, item.Identifier.MessageId, item.Identifier.SenderId)
		for i := 0; i < len(item.LastHop); i++ {
			fmt.Printf("%s | %s\n", item.LastHop[i].Id, item.LastHop[i].NodeType)
		}
		fmt.Println("----------------------------")
	}
}

func (b *brokerService) Broker() {
	broker := b.brokerUsercase.GetBroker()
	for _, broker := range broker.Brokers {
		fmt.Println(broker.Id, broker.Ip, broker.Port)
	}
}

func (b *brokerService) Publisher() {
	broker := b.brokerUsercase.GetBroker()
	for _, publisher := range broker.Publishers {
		fmt.Println(publisher.Id, publisher.Ip, publisher.Port)
	}
}

func (b *brokerService) Subscriber() {
	broker := b.brokerUsercase.GetBroker()
	for _, subscriber := range broker.Subscribers {
		fmt.Println(subscriber.Id, subscriber.Ip, subscriber.Port)
	}
}

func (b *brokerService) GetBroker() *model.Broker {
	return b.brokerUsercase.GetBroker()
}
