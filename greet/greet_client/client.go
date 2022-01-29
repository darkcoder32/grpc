package main

import (
	"context"
	"fmt"
	"grpc/greet/greetpb"
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
