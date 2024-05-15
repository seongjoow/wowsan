package main

import (
	"context"
	"fmt"
	"log"
	"time"
	pb "wowsan/test_chan/client_go/protogo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const GO_SERVER_PORT = "50051"

func rpcConnectTo(ip string) (pb.GoTestServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {
	// connect to the rpc server
	conn, err := grpc.Dial(
		ip,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// connect error handling
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, nil, nil, nil, err
	}

	// new the clinet (client connect interface)
	c := pb.NewGoTestServiceClient(conn)

	// timeout setting
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	return c, conn, ctx, cancel, err
}

func main() {
	// rpc server address
	addr := "localhost" + ":" + GO_SERVER_PORT

	// connect to the rpc server
	client, conn, ctx, cancel, err := rpcConnectTo(addr)
	if err != nil {
		log.Printf("connect error: %v", err)
	}

	/*
		defer:
			after the main function is executed,
			the function in defer will be executed in reverse order
	*/
	defer conn.Close() // close the connection when the main function is executed
	defer cancel()     // cancel the context when the main function is executed

	// call the GetTest method
	r, err := client.GetTestGO(
		ctx,
		&pb.TestRequest{
			TestName: "test go",
		},
	)
	if err != nil {
		log.Printf("grpc call error: %v", err)
	}
	log.Printf("Response: %s", r)
}
