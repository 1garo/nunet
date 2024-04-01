package main

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/1garo/nunet/config"
	"github.com/1garo/nunet/pb"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func server(address string) (pb.DeployerClient, func()) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	baseServer := grpc.NewServer()
	pb.RegisterDeployerServer(baseServer, NewServer())
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		baseServer.Stop()
	}
	c := pb.NewDeployerClient(conn)

	return c, closer
}

func TestDeployJob(t *testing.T) {
	err := config.NewConfig("../.env")
	if err != nil {
		t.Fatal(err)
	}

	c, closer := server("localhost:50051")
	defer closer()
	_, closer = server("localhost:50052")
	defer closer()

	testCases := []struct {
		input    *pb.JobRequest
		deployed bool
		code codes.Code
	}{
		{
			input: &pb.JobRequest{
				Arguments:   []string{"ale"},
				ProgramName: "echo",
				Replicated:  false,
			},
			deployed: true,
		},
		{
			input: &pb.JobRequest{
				Arguments:   []string{"l"},
				ProgramName: "ls",
				Replicated:  false,
			},
			deployed: false,
			code: codes.InvalidArgument,
		},
	}

	for _, tt := range testCases {
		t.Run("test", func(t *testing.T) {
			ctx := context.Background()
			r, err := c.DeployJob(ctx, tt.input)

			if err != nil {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, tt.code, e.Code())
				} else {
					t.Errorf("not able to parse error returned %v", err)
				}
			} else {
				assert.Equal(t, tt.deployed, r.Deployed)
			}
		})
	}
}
