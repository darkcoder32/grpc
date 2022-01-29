package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello world")
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doClientStream(c)
	// doServerStream(c)
	// doUnary(c)
}

func doClientStream(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("failed to init stream: %v\n", err)
		return
	}
	for i := 0; i < 10; i++ {
		req := greetpb.LongGreetRequest{
			MagicNumber: int32(i),
		}

		err = stream.Send(&req)
		if err != nil {
			log.Fatalf("failed to send stream to server: %v\n", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to listen to the server: %v\n", err)
	}
	fmt.Printf("Response from ther server; %v", res.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sumit",
			LastName:  "Kumar",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("GreetManyTimes call failed: %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream; %v\n", err)
		}
		log.Printf("Response from the server: %v\n", msg)
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sumit",
			LastName:  "Kumar",
		},
	}

	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("Client Error in calling Greet: %v\n", err)
	}
	log.Printf("Response from server: %v\n", res.Result)
}
