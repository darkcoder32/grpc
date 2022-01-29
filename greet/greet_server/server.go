package main

import (
	"fmt"
	"grpc/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &greetpb.UnimplementedGreetServiceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}
