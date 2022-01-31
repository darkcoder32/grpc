package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc/chat-service/chatpb"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var count int32

var mu sync.Mutex

func getId() int32 {
	mu.Lock()
	count++
	mu.Unlock()
	return count
}

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

	name := os.Args[1]
	user := &chatpb.User{
		Id:   getId(),
		Name: name,
	}
	stream, err := c.CreateStream(context.Background(), &chatpb.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v\n", stream)
	}
	go func(str chatpb.ChatService_CreateStreamClient) {
		defer func() {
			ch1 <- true
		}()
		for {
			msg, streamErr := str.Recv()
			if streamErr != nil {
				fmt.Printf("Failed to receive message: %v\n", err)
				break
			}
			// msg
			fmt.Printf("Message received from %v: %v\n", msg.GetUser().GetName(), msg.GetText())
		}
	}(stream)
	go func() {
		defer func() {
			ch2 <- true
		}()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &chatpb.Message{
				User:      user,
				Text:      scanner.Text(),
				Timestamp: time.Now().String(),
			}
			_, err := c.BroadCast(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error in broadcastin: %v\n", err)
				break
			}
		}
	}()
	<-ch1
	<-ch2
}
