package main

import (
	pb "GRPCClientServer/gen/proto"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type testApiServer struct {
	pb.UnsafeTestApiServer
}

// A Response struct to map the Entire Response
type Response struct {
	Content []string `json:"content"`
	Found   string   `json:"isFound"`
}

func (s *testApiServer) FindLog(ctx context.Context, req *pb.LambdaRequest) (*pb.LambdaResponse, error) {

	url := fmt.Sprintf("https://imh5ufzsd9.execute-api.us-east-1.amazonaws.com/prod/checkifpresent?T=%s&dT=%s", req.Time, req.Deltatime)

	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println(response.Body)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(responseData)

	var responseObject Response

	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.Content)

	if responseObject.Found == "false" {
		var emptyArray []string

		responseBody := Response{
			Content: emptyArray,
			Found:   "False",
		}

		responseToClient, _ := json.Marshal(responseBody)

		return &pb.LambdaResponse{Result: string(responseToClient)}, nil
	} else {
		responseBody := Response{
			Content: responseObject.Content,
			Found:   responseObject.Found,
		}

		responseToClient, _ := json.Marshal(responseBody)

		return &pb.LambdaResponse{Result: string(responseToClient)}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTestApiServer(grpcServer, &testApiServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
}
