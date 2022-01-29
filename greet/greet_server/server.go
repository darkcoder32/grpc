package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
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

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}
