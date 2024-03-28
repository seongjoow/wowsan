package usecase

import (
	"fmt"
	"wowsan/constants"
	_brokerClient "wowsan/pkg/broker/grpc/client"
	model "wowsan/pkg/model"

	uuid "github.com/satori/go.uuid"
)

type SubscriberUsecase interface {
	Sub(subject string, operator string, value string, brokerIp string, brokerPort string)
	ReceivePublication(id string, subject string, operator string, value string) error
}

type subscriberUsecase struct {
	subscriber   *model.Subscriber
	brokerClient _brokerClient.BrokerClient
}

func NewSubscriberUsecase(
	subscriber *model.Subscriber,
	brokerClient _brokerClient.BrokerClient,
) SubscriberUsecase {
	return &subscriberUsecase{
		subscriber:   subscriber,
		brokerClient: brokerClient,
	}
}

func (uc *subscriberUsecase) Sub(subject string, operator string, value string, brokerIp string, brokerPort string) {
	subscriber := uc.subscriber

	if subscriber.Broker == nil {
		response, err := uc.brokerClient.RPCAddSubscriber(brokerIp, brokerPort, subscriber.Id, subscriber.Ip, subscriber.Port)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		broker := model.NewBroker(response.Id, response.Ip, response.Port)
		subscriber.SetBroker(broker)
		fmt.Printf("Set broker: %s %s %s\n", response.Id, response.Ip, response.Port)
	}

	_, err := uc.brokerClient.RPCSendSubscription(
		subscriber.Broker.Ip,
		subscriber.Broker.Port,
		subscriber.Id,
		subscriber.Ip,
		subscriber.Port,
		subject,
		operator,
		value,
		constants.SUBSCRIBER,
		uuid.NewV4().String(),
		subscriber.Id,
	)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func (uc *subscriberUsecase) ReceivePublication(id string, subject string, operator string, value string) error {
	fmt.Println("NOTIFICATION: Received a publication")
	return nil
}
