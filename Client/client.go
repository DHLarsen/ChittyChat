package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	gRPC "github.com/DHLarsen/ChittyChat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client gRPC.ModelClient
var ServerConn *grpc.ClientConn

var clientName string

func ConnectToServer() {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(":8888", opts...)
	if err != nil {
		print(err)
	}

	client = gRPC.NewModelClient(conn)
	ServerConn = conn

	log.Println("the connection is: ", conn.GetState().String())

	go updateListen()
}

func updateListen() {
	updateRequest := gRPC.UpdateRequest{}

	stream, err := client.GetUpdate(context.Background())
	if err != nil {
		log.Fatal(err) // dont use panic in your real project
	}

	stream.Send(&updateRequest)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			//panic(err)
		}
		if resp != nil{
			printOutput(resp)
		}
	}
}

func main() {
	startup()
	ConnectToServer()

	defer ServerConn.Close()
	reader := bufio.NewReader(os.Stdin)

	stream, err := client.SendMessage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//Infinite loop to listen for clients input.
	for {
		prepareInput()

		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		stream.Send(&gRPC.Message{ClientName: clientName, Message: input})
		continue
	}
}

func prepareInput() {
	fmt.Print("-> ")
}

func printOutput(msg *gRPC.Message) {
	name := msg.ClientName
	if name == clientName {
		name = "You"
	}
	fmt.Println(name, ": ", msg.Message)
	prepareInput()
}

func startup() {
	fmt.Print("Please enter your name: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	clientName = strings.TrimSpace(input.Text())
}
