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
	"google.golang.org/grpc/credentials/insecure"
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
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
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
		input *pb.JobRequest
		//expected error
		deployed bool
	}{
		{
			&pb.JobRequest{
				Arguments:   []string{"ale"},
				ProgramName: "echo",
				Replicated:  false,
			}, 
			//nil,
			true,
		},
		// TODO: test the error messages
		//{
		//	&pb.JobRequest{
		//		Arguments:   []string{"l"},
		//		ProgramName: "ls",
		//		Replicated:  false,
		//	}, 
		//	//errors.New("could not run job: ls: cannot access 'l': No such file or directory"),
		//	false,
		//},
	}

	for _, tt := range testCases {
		ctx := context.Background()
		r, _ := c.DeployJob(ctx, tt.input)

		//assert.Equal(t, tt.expected, err)
		assert.Equal(t, tt.deployed, r.Deployed)
	}
}
