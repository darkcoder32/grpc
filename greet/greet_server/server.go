package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetingResponse, error) {
	fmt.Printf("Greet function is invoked: %v\n", req)
	res := &greetpb.GreetingResponse{
		Result: "Hello" + req.GetGreeting().FirstName,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function is invoked; %v\n", req)
	for i := 0; i < 10; i++ {
		res := greetpb.GreetingManyTimesResponse{
			Result: "Hello " + req.GetGreeting().GetFirstName() + " number " + strconv.Itoa(i),
		}
		err := stream.Send(&res)
		if err != nil {
			log.Fatalf("Stream Send failed: %v\n", err)
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function is invoked; %v\n", stream)
	var sum int32
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: sum,
			})
		}
		if err != nil {
			log.Fatalf("Stream Reading error: %v\n", err)
		}
		sum += msg.GetMagicNumber()
	}
}

func (*server) BiGreet(stream greetpb.GreetService_BiGreetServer) error {
	fmt.Printf("BiGreet function is invoked: %v\n", stream)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Stream Reading error: %v\n", err)
		}
		fmt.Printf("Small number from client received: %v\n", req.GetSmallNumber())
		err = stream.Send(&greetpb.BiGreetResponse{
			LargeNumber: req.GetSmallNumber() + 1000,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function is invoked: %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Negative number received from the client: %v\n", number))
	}
	return &greetpb.SquareRootResponse{Result: math.Sqrt(float64(number))}, nil
}

func (*server) GreetDeadline(ctx context.Context, req *greetpb.GreetDeadlineRequest) (*greetpb.GreetDeadlineResponse, error) {
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Printf("client canceled the request\n")
			return nil, status.Errorf(codes.DeadlineExceeded, "Deadline exceeded")
		}
		time.Sleep(time.Second)
	}
	return &greetpb.GreetDeadlineResponse{Result: "Hello " + req.GetGreeting().GetFirstName()}, nil
}

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed to load certificates: %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)
	s := grpc.NewServer(opts)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}
