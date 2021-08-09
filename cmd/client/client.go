package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/codeedu/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// Unary
	//AddUser(client)

	// Server Stream
	//AddUserVerbose(client)

	// Client Stream
	AddUsers(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Ethan",
		Email: "ethan@ethan.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Ethan",
		Email: "ethan@ethan.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	// Loop status

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg %v", err)
		}

		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "u1",
			Name:  "Ed1",
			Email: "ed@ed1.com",
		},
		&pb.User{
			Id:    "u2",
			Name:  "Ed2",
			Email: "ed@ed2.com",
		},
		&pb.User{
			Id:    "u3",
			Name:  "Ed3",
			Email: "ed@ed3.com",
		},
		&pb.User{
			Id:    "u4",
			Name:  "Ed4",
			Email: "ed@ed4.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Eror receiving request: %v", err)
	}

	fmt.Println(res)
}
