package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/1garo/nunet/config"
	pb "github.com/1garo/nunet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Define the gRPC jobServer
type jobServer struct {
	pb.UnimplementedDeployerServer
}

func NewServer() pb.DeployerServer {
	j := &jobServer{}
	return j
}

// Replication Logic in both containers
func (s *jobServer) DeployJob(ctx context.Context, req *pb.JobRequest) (*pb.JobResponse, error) {
	fmt.Printf("target port: %s\n", config.Cfg.TargetAddr)
	fmt.Printf("programName(%s): args (%+v)\n", req.ProgramName, req.Arguments)
	var stderr bytes.Buffer

	cmd := exec.Command(req.ProgramName, req.Arguments...)

	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("could not run job: %s\n", stderr.String())
		return &pb.JobResponse{Deployed: false}, errors.New(stderr.String())
	}

	fmt.Println("program output: (start)")
	fmt.Print(stderr.String())
	fmt.Println("program output: (end)")
	if !req.Replicated {
		req.Replicated = true
		return deployJobOnOtherContainer(req, config.Cfg.TargetAddr)
	} else {
		return &pb.JobResponse{Deployed: true}, nil
	}
}

// Implement the gRPC client
func deployJobOnOtherContainer(req *pb.JobRequest, targetAddr string) (*pb.JobResponse, error) {
	c := grpc.WithTransportCredentials(insecure.NewCredentials())
	target := fmt.Sprintf("%s:%s", config.Cfg.ClientName, targetAddr)
	conn, err := grpc.Dial(target, c)
	if err != nil {
		return &pb.JobResponse{Deployed: false}, err
	}
	defer conn.Close()

	client := pb.NewDeployerClient(conn)
	// In this case we don't care to handle errors, we are just propagating the job
	return client.DeployJob(context.Background(), req)
}
