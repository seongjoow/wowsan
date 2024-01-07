package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "wowsan/test_chan/server_go/protogo"

	"google.golang.org/grpc"
)

const GO_SERVER_PORT = "50051"

type server struct {
	*pb.UnimplementedGoTestServiceServer
}

func (s *server) GetTestGO(ctx context.Context, reqeust *pb.TestRequest) (*pb.TestResponse, error) {
	// for test
	fmt.Printf("Log on rpc server: %v\n", reqeust.TestName)
	// response for client
	res := &pb.TestResponse{
		Message: "Test Name: " + reqeust.TestName,
	}
	return res, nil
}

func initServer() {
	lis, err := net.Listen("tcp", ":"+GO_SERVER_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("listen on 50051")

	s := grpc.NewServer()
	pb.RegisterGoTestServiceServer(s, &server{})
	s.Serve(lis) // start the server (blocking function)
}

func main() {
	// it's a blocking function, so the following code will not be executed
	initServer()

	// go initServer()
	// fmt.Println("rpc server is running...")
}
