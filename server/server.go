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

var vTime = []int64{0}
var vTimeIndex = 0

var updateChansMutex sync.Mutex
var updateChans = []chan *gRPC.Message{}

//var messages = []*gRPC.Message{}

func (s *Server) SendMessage(msgStream gRPC.Model_SendMessageServer) error {
	clientName := ""
	for {
		msg, err := msgStream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(clientName, " disconnected")
			return err
		} else if clientName == "" {
			clientName = msg.ClientName
		}

		updateVTime(msg.Time)
		//increment vector timestamp when recieving message

		updateChansMutex.Lock()
		incrementVTime()
		updateChansMutex.Unlock()
		msg.Time = vTime

		broadcastMsg(msg)
		log.Printf("Received message from %s: %s", msg.ClientName, msg.Message)
	}

	return nil
}
func (s *Server) GetUpdate(updateStream gRPC.Model_GetUpdateServer) error {
	msg, err := updateStream.Recv()
	if err == io.EOF {
		log.Fatal("Error in creating connection")
	}
	updateChan := make(chan *gRPC.Message)
	updateChans = append(updateChans, updateChan)
	syncMsg := gRPC.Message{
		ClientName: "",
		Message:    "Participant " + msg.ClientName + " joined the chat room!",
		Time:       vTime,
	}
	broadcastMsg(&syncMsg)
	for {
		var msg = <-updateChan
		sendMessage(msg, updateStream)
	}
}

func broadcastMsg(msg *gRPC.Message) {
	for _, updateChan := range updateChans {
		updateChan <- msg
	}
}

func sendMessage(msg *gRPC.Message, updateStream gRPC.Model_GetUpdateServer) {
	updateChansMutex.Lock()
	incrementVTime()
	msg.Time = vTime
	log.Println("Sending: ", msg.Message)
	updateStream.Send(msg)
	updateChansMutex.Unlock()
}

func updateVTime(newVTime []int64) {
	updateChansMutex.Lock()
	for len(vTime) < len(newVTime) {
		vTime = append(vTime, 0)
	}
	for i, time := range newVTime {
		if time > vTime[i] {
			vTime[i] = time
		}
	}
	log.Println("Set time to: ", vTime)
	updateChansMutex.Unlock()
}

func incrementVTime() {
	vTime[vTimeIndex]++
	log.Println("Set time to: ", vTime)
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
