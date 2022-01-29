package main

import (
	"context"
	"fmt"
	"grpc/calculator/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client Hello")
	conn, err := grpc.Dial("0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to connect to client: %v", err)
	}
	c := pb.NewCalServiceClient(conn)

	req := pb.Request{
		FirstNumber:  3,
		SecondNumber: 10,
	}
	res, err := c.Calculate(context.Background(), &req)
	if err != nil {
		log.Fatalf("Calculate called failed: %v\n", err)
	}
	fmt.Printf("Response from client: %v\n", res)
}
