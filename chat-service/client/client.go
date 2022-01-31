package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc/chat-service/chatpb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client Hello")
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect the server")
	}
	defer conn.Close()
	c := chatpb.NewChatServiceClient(conn)

	ch1 := make(chan bool)
	ch2 := make(chan bool)

	go func() {
		defer func() {
			ch1 <- true
		}()
		stream, err := c.CreateStream(context.Background(), &chatpb.Connect{
			User: &chatpb.User{
				Id:   1,
				Name: "Sumit Kumar",
			},
			Active: true,
		})
		if err != nil {
			log.Fatalf("Failed to call create stream: %v\n", err)
			return
		}
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive message: %v\n", err)
				break
			}
			fmt.Printf("Message received from id %d: %v", msg.GetUser().GetId(), msg.GetText())
		}
	}()
	go func() {
		defer func() {
			ch2 <- true
		}()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &chatpb.Message{
				User: &chatpb.User{
					Id:   1,
					Name: "Sumit Kumar",
				},
				Text:      scanner.Text(),
				Timestamp: time.Now().String(),
			}
			_, err := c.BroadCast(context.Background(), msg)
			if err != nil {
				log.Fatalf("Error in broadcastin: %v\n", err)
				break
			}
		}
	}()
	<-ch1
	<-ch2
}
