package main

import (
	"context"
	"fmt"
	"grpc/chat-service/chatpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Connection struct {
	stream chatpb.ChatService_CreateStreamServer
	id     int32
	active bool
	err    chan error
}
type Server struct {
	Connection []*Connection
}

func (s *Server) CreateStream(con *chatpb.Connect, stream chatpb.ChatService_CreateStreamServer) error {
	fmt.Println("CreatStream invoked")
	conn := &Connection{
		stream: stream,
		id:     con.GetUser().GetId(),
		active: true,
		err:    make(chan error),
	}
	s.Connection = append(s.Connection, conn)
	return nil
}

func (s *Server) BroadCast(ctx context.Context, msg *chatpb.Message) (*chatpb.Close, error) {
	fmt.Println("Brodcast invoked")
	for _, conn := range s.Connection {
		if conn.active {
			err := conn.stream.Send(msg)
			if err != nil {
				conn.active = false
				conn.err <- err
				log.Fatalf("Error sending messgae to user: user id %d: %v\n", conn.id, err)
			}
		}
	}
	return &chatpb.Close{}, nil
}

func main() {
	fmt.Println("Server Hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen to the port: %v\n", err)
	}
	grpc := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpc, &Server{})
	if err := grpc.Serve(lis); err != nil {
		log.Fatalf("Failed to serve the server: %v\n", err)
	}
}
