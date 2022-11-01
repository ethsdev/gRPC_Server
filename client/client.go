package main

import (
	pb "GRPCClientServer/gen/proto"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

// A Response struct to map the Entire Response
type Response struct {
	Content []string `json:"content"`
	Found   string   `json:"isFound"`
}

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}

	client := pb.NewTestApiClient(conn)

	fmt.Println("Enter the searching time T: ")

	// var then variable name then variable type
	var first string

	// Taking input from user
	fmt.Scanln(&first)

	fmt.Println("Enter the delta time dT: ")

	// var then variable name then variable type
	var second string

	// Taking input from user
	fmt.Scanln(&second)

	resp, err := client.FindLog(context.Background(), &pb.LambdaRequest{Time: first, Deltatime: second})

	var responseObject Response

	json.Unmarshal([]byte(resp.Result), &responseObject)

	if err != nil {
		log.Println(err)
	}

	if responseObject.Found == "True" {
		fmt.Println("Congrats! Match was found! Here are hashes: ")
		fmt.Println(responseObject.Content)
	} else {
		fmt.Println("Unfortunately no hashes were found: ")
		fmt.Println(responseObject.Content)
	}
}
