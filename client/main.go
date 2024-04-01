package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	pb "github.com/1garo/nunet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr        = flag.String("addr", "localhost:50051", "the address to connect to")
	programName = flag.String("program", "", "the program to run")
	arguments   = flag.String("args", "", "the arguments of the program (e.g -args=arg1,arg2)")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDeployerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.DeployJob(ctx, &pb.JobRequest{
		Arguments:   strings.Split(*arguments, ","),
		ProgramName: *programName,
		Replicated:  false,
	})

	if err != nil {
		log.Fatalf("could not get deployed status: %v", err)
	}
	log.Printf("Deployed: %t\n", r.Deployed)
}
