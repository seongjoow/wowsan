package service

import (
	"fmt"
	"log"
	"net"

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
	PrtWithSubject(subject string)
	SrtWithSubject(subject string)
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

	RPCErrorLogger, err := logger.NewLogger("RPCerror", "./log/RPCErrorLogger/"+baseDir)
	if err != nil {
		panic(err)
	}

	brokerClient := _brokerClient.NewBrokerClientWithLogger(RPCErrorLogger)
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

	// go func() {
	// 	for {
	// 		b := brokerUsecase.GetBroker()
	// 		if time.Since(b.Close) > 30*time.Second {
	// 			// terminate the program
	// 			fmt.Println("Broker is closed")
	// 			http.PostForm("http://localhost:8080/done", map[string][]string{"port": {port}})
	// 			// os.Exit(1)
	// 		}
	// 	}
	// }()
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
	broker.SRT.PrintSRT()
	// // Define headers and lengths
	// columnHeaders := []string{"Subject", "Operator", "Value", "MessageId", "SenderId", "[]LastHop"}
	// columnLengths := []int{15, 10, 10, 15, 10, 25}
	// // Print the table
	// table := utils.NewTable(columnHeaders, columnLengths)
	// table.SetTitle("SRT")
	// for _, item := range broker.SRT {
	// 	var lastHop string
	// 	for i, hop := range item.LastHop {
	// 		if i > 0 {
	// 			lastHop += ", "
	// 		}
	// 		lastHop += fmt.Sprintf("%s(%s)", hop.Id, hop.NodeType)
	// 	}
	// 	row := []string{
	// 		item.Advertisement.Subject,
	// 		item.Advertisement.Operator,
	// 		item.Advertisement.Value,
	// 		item.Identifier.MessageId,
	// 		item.Identifier.SenderId,
	// 		lastHop,
	// 	}
	// 	table.AddRow(row)
	// }

	// // Print the table
	// table.PrintTable()

}

func (b *brokerService) SrtWithSubject(subject string) {
	broker := b.brokerUsercase.GetBroker()
	broker.SRT.PrintSRTWhereSubject(subject)
}

func (b *brokerService) Prt() {
	broker := b.brokerUsercase.GetBroker()
	broker.PRT.PrintPRT()
}

func (b *brokerService) PrtWithSubject(subject string) {
	broker := b.brokerUsercase.GetBroker()
	broker.PRT.PrintPRTWhereSubject(subject)
}

func (b *brokerService) Broker() {
	broker := b.brokerUsercase.GetBroker()
	for _, broker := range broker.Brokers {
		fmt.Println(broker.Id, broker.Ip, broker.Port)
	}
}

func (b *brokerService) Publisher() {
	broker := b.brokerUsercase.GetBroker()
	broker.PrintPublisher()
	// // for _, publisher := range broker.Publishers {
	// // 	fmt.Println(publisher.Id, publisher.Ip, publisher.Port)
	// // }
	// fmt.Println("[Publisher]")
	// // Define headers and lengths
	// columnHeaders := []string{"Id", "Ip", "Port"}
	// columnLengths := []int{15, 15, 15}
	// // Print the table
	// table := utils.NewTable(columnHeaders, columnLengths)
	// for _, publisher := range broker.Publishers {
	// 	row := []string{
	// 		publisher.Id,
	// 		publisher.Ip,
	// 		publisher.Port,
	// 	}
	// 	table.AddRow(row)
	// }
	// table.PrintTable()
}

func (b *brokerService) Subscriber() {
	broker := b.brokerUsercase.GetBroker()
	broker.PrintSubscriber()
	// // for _, subscriber := range broker.Subscribers {
	// // 	fmt.Println(subscriber.Id, subscriber.Ip, subscriber.Port)
	// // }
	// fmt.Println("[Subscriber]")
	// // Define headers and lengths
	// columnHeaders := []string{"Id", "Ip", "Port"}
	// columnLengths := []int{15, 15, 15}
	// // Print the table
	// table := utils.NewTable(columnHeaders, columnLengths)
	// for _, subscriber := range broker.Subscribers {
	// 	row := []string{
	// 		subscriber.Id,
	// 		subscriber.Ip,
	// 		subscriber.Port,
	// 	}
	// 	table.AddRow(row)
	// }
	// table.PrintTable()
}

func (b *brokerService) GetBroker() *model.Broker {
	return b.brokerUsercase.GetBroker()
}
