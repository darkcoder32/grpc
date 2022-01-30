package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello world")
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	// successful run
	doUnaryDeadline(c, 5*time.Second)
	// deadline exceeeded
	doUnaryDeadline(c, 1*time.Second)
	// doSquareRoot(c)
	// doBiStream(c)
	// doClientStream(c)
	// doServerStream(c)
	// doUnary(c)
}

func doUnaryDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := greetpb.GreetDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sumit",
			LastName:  "Kumar",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetDeadline(ctx, &req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Printf("Deadline exceeded by server\n")
			} else {
				fmt.Printf("unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Client Error in calling Greet: %v\n", err)
		}
	}
	log.Printf("Response from server: %v\n", res.GetResult())
}

func doSquareRoot(c greetpb.GreetServiceClient) {
	req := &greetpb.SquareRootRequest{Number: -1}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		fromErr, ok := status.FromError(err)
		if ok {
			fmt.Println(fromErr.Message())
			fmt.Println(fromErr.Code())
			if fromErr.Code() == codes.InvalidArgument {
				fmt.Printf("client probably sent a negative number\n")
			}
		} else {
			fmt.Printf("SquareRoot calling failed for server: %v\n", err)
		}
	}
	fmt.Printf("Response received from client: %v\n", res.GetResult())
}

func doBiStream(c greetpb.GreetServiceClient) {
	stream, err := c.BiGreet(context.Background())
	if err != nil {
		log.Fatalf("failed to init stream: %v\n", err)
		return
	}

	ch1 := make(chan bool)
	ch2 := make(chan bool)

	go BiSend(stream, ch1)
	go BiRecv(stream, ch2)
	<-ch1
	<-ch2
}

func BiSend(stream greetpb.GreetService_BiGreetClient, ch chan bool) error {
	defer func() {
		ch <- true
	}()
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

func BiRecv(stream greetpb.GreetService_BiGreetClient, ch chan bool) error {
	defer func() {
		ch <- true
	}()
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
