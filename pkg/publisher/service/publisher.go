package service

import (
	_brokerClient "wowsan/pkg/broker/grpc/client"
	"wowsan/pkg/model"
	_publisherUsecase "wowsan/pkg/publisher/usecase"
)

type PublisherService struct {
	Publisher        *model.Publisher
	PublisherUsecase _publisherUsecase.PublisherUsecase
	BrokerClient     _brokerClient.BrokerClient
}

func NewPublisherService(
	ip string,
	port string,
) *PublisherService {
	id := ip + ":" + port
	publisher := model.NewPublisher(id, ip, port)

	brokerClient := _brokerClient.NewBrokerClient()

	publisherUsecase := _publisherUsecase.NewPublisherUsecase(
		publisher,
		brokerClient,
	)

	return &PublisherService{
		Publisher:        publisher,
		PublisherUsecase: publisherUsecase,
		BrokerClient:     brokerClient,
	}
}
