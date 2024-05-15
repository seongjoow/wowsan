package main

import (
	"context"
	"log"
	"time"

	pb "test/proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func subscribe(client pb.BrokerServiceClient, conditions map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.Subscribe(ctx, &pb.Subscription{Conditions: conditions})
	if err != nil {
		log.Printf("subscribe error: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("error receiving message: %v", err)
		}
		log.Printf("Received message: %s", msg.Content)
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerServiceClient(conn)

	// 조건에 맞는 메시지 구독
	conditions := map[string]string{"type": "news"}
	go subscribe(client, conditions)

	// 다른 작업 수행...
}
