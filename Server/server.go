package main

import (
	"io"
	"log"
	"net"
	"sync"

	gRPC "github.com/DHLarsen/ChittyChat/proto"
	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedModelServer

	name  string
	port  string
	mutex sync.Mutex
}

var messages = []*gRPC.Message{}

func (s *Server) SendMessage(msgStream gRPC.Model_SendMessageServer) error {
	for {
		msg, err := msgStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v", err)
			return err
		}

		messages = append(messages, msg)
		log.Printf("Received message from %s: %s", msg.ClientName, msg.Message)
	}

	return nil
}
func (s *Server) GetUpdate(updateStream gRPC.Model_GetUpdateServer) error {
	for _, msg := range messages {
		if err := updateStream.Send(msg); err != nil {
			log.Println("Error: ", err)
			return err
		}
	}

	return nil
}

func launchServer() {
	list, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := &Server{
		name: "Server",
		port: "8888",
	}
	gRPC.RegisterModelServer(grpcServer, server)
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func main() {
	launchServer()
}
