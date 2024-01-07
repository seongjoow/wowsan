package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "test/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedBrokerServiceServer
	subscribers sync.Map // 구독자 관리
}

// 구독자 정보
type subscriber struct {
	conditions map[string]string
	stream     pb.BrokerService_SubscribeServer
}

// 메시지가 조건과 일치하는지 확인
func matchConditions(attributes, conditions map[string]string) bool {
	for key, value := range conditions {
		if val, ok := attributes[key]; !ok || val != value {
			return false
		}
	}
	return true
}

func (s *server) Publish(ctx context.Context, in *pb.Publication) (*pb.PublishAck, error) {
	log.Printf("Received publication: %v", in)
	s.subscribers.Range(func(_, value interface{}) bool {
		sub := value.(*subscriber)
		if matchConditions(in.Attributes, sub.conditions) {
			sub.stream.Send(in)
		}
		return true
	})
	return &pb.PublishAck{Success: true}, nil
}

func (s *server) Subscribe(in *pb.Subscription, stream pb.BrokerService_SubscribeServer) error {
	sub := &subscriber{
		conditions: in.Conditions,
		stream:     stream,
	}
	subID := // 적절한 구독자 ID 생성
		s.subscribers.Store(subID, sub)
	// 연결 종료 시 구독자 삭제
	defer s.subscribers.Delete(subID)

	// 구독 유지 (클라이언트 연결 종료까지 대기)
	<-stream.Context().Done()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBrokerServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
