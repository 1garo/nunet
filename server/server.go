package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/1garo/nunet/pb" // Import your protobuf generated code
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Define the gRPC server
type server struct {
	pb.UnimplementedDeployerServer
}

// Replication Logic in both containers
func (s *server) DeployJob(ctx context.Context, req *pb.JobRequest) (*pb.JobResponse, error) {
	targetAddr := os.Getenv("OTHER_CONTAINER_ADDRESS")
	if targetAddr == "" {
		return nil, fmt.Errorf("OTHER_CONTAINER_ADDRESS environment variable not set")
	}

	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		log.Fatal("GRPC_PORT environment variable not set")
	}
	fmt.Printf("port: %s\n", portStr)
	fmt.Printf("processing request (%s): %+v\n", req.ProgramName, req.Arguments)
	if !req.Replicated {
		req.Replicated = true
		return deployJobOnOtherContainer(req, targetAddr)
	} else {
		return &pb.JobResponse{Status: "200"}, nil
	}
}

var CLIENT_NAME string

func main() {
	//env := os.Getenv("ENV")
	//if env == "" {
	//	log.Fatal("ENV environment variable not set")
	//}
	//err := godotenv.Load(env)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		log.Fatal("GRPC_PORT environment variable not set")
	}

	CLIENT_NAME = os.Getenv("CLIENT_NAME")
	if CLIENT_NAME == "" {
		log.Fatal("CLIENT_NAME environment variable not set")
	}

	lis, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDeployerServer(s, &server{})
	log.Printf("gRPC server started on port %s", portStr)
	s.Serve(lis)
}

// Implement the gRPC client
func deployJobOnOtherContainer(req *pb.JobRequest, targetAddr string) (*pb.JobResponse, error) {
	c := grpc.WithTransportCredentials(insecure.NewCredentials())
	p := fmt.Sprintf("%s:%s", CLIENT_NAME, targetAddr)
	conn, err := grpc.Dial(p, c)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDeployerClient(conn)
	return client.DeployJob(context.Background(), req)
}
