package usecase

import (
	"fmt"
	"wowsan/constants"
	model "wowsan/pkg/model"

	uuid "github.com/satori/go.uuid"

	_brokerClient "wowsan/pkg/broker/grpc/client"
)

type PublisherUsecase interface {
	Adv(subject string, operator string, value string, brokerIp string, brokerPort string)
	Pub(subject string, operator string, value string, brokerIp string, brokerPort string)
}

type publisherUsecase struct {
	publisher    *model.Publisher
	brokerClient _brokerClient.BrokerClient
}

func NewPublisherUsecase(
	publisher *model.Publisher,
	brokerClient _brokerClient.BrokerClient,
) PublisherUsecase {
	return &publisherUsecase{
		publisher:    publisher,
		brokerClient: brokerClient,
	}
}

func (uc *publisherUsecase) Adv(subject string, operator string, value string, brokerIp string, brokerPort string) {
	publisher := uc.publisher

	if publisher.Broker == nil {
		response, err := uc.brokerClient.RPCAddPublisher(brokerIp, brokerPort, publisher.Id, publisher.Ip, publisher.Port)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		broker := model.NewBroker(response.Id, response.Ip, response.Port)
		publisher.SetBroker(broker)
		fmt.Printf("Set broker: %s %s %s\n", response.Id, response.Ip, response.Port)
	}

	hopCount := int64(0)
	_, err := uc.brokerClient.RPCSendAdvertisement(
		publisher.Broker.Ip,
		publisher.Broker.Port,
		publisher.Id,
		publisher.Ip,
		publisher.Port,
		subject,
		operator,
		value,
		constants.PUBLISHER,
		hopCount,
		uuid.NewV4().String(), // TODO: messageId
		publisher.Id,
	)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func (uc *publisherUsecase) Pub(subject string, operator string, value string, brokerIp string, brokerPort string) {
	publisher := uc.publisher

	if publisher.Broker == nil {
		fmt.Printf("This publisher isn't registered to a broker.\n")
	}
	if publisher.Broker.Ip != brokerIp || publisher.Broker.Port != brokerPort {
		fmt.Printf("This publisher isn't registered to the broker.\n")
	}

	_, err := uc.brokerClient.RPCSendPublication(
		publisher.Broker.Ip,
		publisher.Broker.Port,
		publisher.Id,
		publisher.Ip,
		publisher.Port,
		subject,
		operator,
		value,
		constants.PUBLISHER,
		uuid.NewV4().String(), // TODO: messageId
	)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}
