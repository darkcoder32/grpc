package main

import (
	"context"
	"fmt"
	"grpc/calculator/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Printf("Calculate called: %v\n", req)
	res := &pb.Response{
		Sum: int64(req.GetFirstNumber()) + int64(req.GetSecondNumber()),
	}
	return res, nil
}

func main() {
	fmt.Println("Server Hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listen failed on server: %v\n", err)
	}
	c := grpc.NewServer()
	pb.RegisterCalServiceServer(c, &server{})
	if err := c.Serve(lis); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
