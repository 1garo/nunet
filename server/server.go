package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/1garo/nunet/config"
	pb "github.com/1garo/nunet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	log.Println("[DeployJob]")
	log.Printf("target port: %s\n", config.Cfg.TargetAddr)
	log.Printf("programName(%s): args (%+v)\n", req.ProgramName, req.Arguments)
	var stderr bytes.Buffer
	var out bytes.Buffer

	cmd := exec.Command(req.ProgramName, req.Arguments...)

	cmd.Stderr = &stderr
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("could not run job: %s\n", stderr.String())
		return nil, status.Errorf(codes.InvalidArgument, stderr.String())
	}

	log.Println("program output: (start)")
	log.Print(out.String())
	log.Println("program output: (end)")
	if !req.Replicated {
		req.Replicated = true
		return deployJobOnOtherContainer(req, config.Cfg.TargetAddr)
	} else {
		return &pb.JobResponse{Deployed: true}, nil
	}
}

// Implement the gRPC client
func deployJobOnOtherContainer(req *pb.JobRequest, targetAddr string) (*pb.JobResponse, error) {
	log.Println("[deployJobOnOtherContainer]")
	c := grpc.WithTransportCredentials(insecure.NewCredentials())
	target := fmt.Sprintf("%s:%s", config.Cfg.ClientName, targetAddr)
	conn, err := grpc.Dial(
		target, 
		c, 
		grpc.FailOnNonTempDialError(true),
        grpc.WithBlock(), 
	)

	if err != nil {
		log.Printf("error while trying to dial client: %s\n", err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	defer conn.Close()

	client := pb.NewDeployerClient(conn)
	// In this case we don't care to handle errors, we are just propagating the job
	return client.DeployJob(context.Background(), req)
}
