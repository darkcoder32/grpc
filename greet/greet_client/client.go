package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
	"io"
	"log"
	"time"

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

	doBiStream(c)
	// doClientStream(c)
	// doServerStream(c)
	// doUnary(c)
}

func doBiStream(c greetpb.GreetServiceClient) {
	stream, err := c.BiGreet(context.Background())
	if err != nil {
		log.Fatalf("failed to init stream: %v\n", err)
		return
	}
	go BiSend(stream)
	go BiRecv(stream)
	time.Sleep(time.Minute)
}

func BiSend(stream greetpb.GreetService_BiGreetClient) error {
	fmt.Printf("BiSend function is invoked; %v\n", stream)
	for i := 0; i < 10; i++ {
		res := greetpb.BiGreetRequest{
			SmallNumber: int32(i + 1),
		}
		err := stream.Send(&res)
		if err != nil {
			log.Fatalf("Stream Send failed: %v\n", err)
			return err
		}
		time.Sleep(time.Second)
	}
	stream.CloseSend()
	return nil
}

func BiRecv(stream greetpb.GreetService_BiGreetClient) error {
	fmt.Printf("BiRecv function is invoked; %v\n", stream)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Stream Reading error: %v\n", err)
		}
		fmt.Printf("Large number from the server: %v\n", msg.GetLargeNumber())
	}
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
