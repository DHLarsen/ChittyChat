package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	gRPC "github.com/DHLarsen/ChittyChat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var server gRPC.ModelClient
var ServerConn *grpc.ClientConn

func ConnectToServer() {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(":8888", opts...)
	if err != nil {
		print(err)
	}

	server = gRPC.NewModelClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())

}

func main() {
	ConnectToServer()
	defer ServerConn.Close()

	reader := bufio.NewReader(os.Stdin)

	//Infinite loop to listen for clients input.
	for {
		fmt.Print("-> ")

		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		if input == "hi" {
			stream, err := server.SendMessage(context.Background())
			if err != nil {
				print(err)
			}

			greet := gRPC.Message{ClientName: "client", Message: "Hello from client!"}

			stream.Send(&greet)
			stream.Send(&greet)
		}
		continue
	}
}
