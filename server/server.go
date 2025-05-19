package main

import (
	pb "GRPCClientServer/gen/proto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// Response represents the JSON structure from the external API
type Response struct {
	Content []string `json:"content"`
	Found   string   `json:"isFound"`
}

type testApiServer struct {
	pb.UnimplementedTestApiServer // use the proper embedded struct
	client *http.Client
}

// NewTestApiServer is a constructor for testApiServer
func NewTestApiServer() *testApiServer {
	return &testApiServer{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// fetchLogFromAPI calls the external API and unmarshals the response
func (s *testApiServer) fetchLogFromAPI(timeStr, deltaTime string) (Response, error) {
	url := fmt.Sprintf("https://imh5ufzsd9.execute-api.us-east-1.amazonaws.com/prod/checkifpresent?T=%s&dT=%s", timeStr, deltaTime)

	resp, err := s.client.Get(url)
	if err != nil {
		return Response{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, errors.New("non-200 response from API")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("error reading response: %w", err)
	}

	var result Response
	if err := json.Unmarshal(data, &result); err != nil {
		return Response{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return result, nil
}

func (s *testApiServer) FindLog(ctx context.Context, req *pb.LambdaRequest) (*pb.LambdaResponse, error) {
	apiResp, err := s.fetchLogFromAPI(req.Time, req.Deltatime)
	if err != nil {
		log.Printf("API fetch error: %v", err)
		return nil, err
	}

	if apiResp.Found == "false" {
		apiResp.Content = []string{}
	}

	resJSON, err := json.Marshal(apiResp)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return nil, err
	}

	return &pb.LambdaResponse{Result: string(resJSON)}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTestApiServer(grpcServer, NewTestApiServer())

	log.Println("gRPC server listening on localhost:9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
