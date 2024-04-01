package main

import (
	"log"
	"net"

	"github.com/1garo/nunet/config"
	pb "github.com/1garo/nunet/pb" // Import your protobuf generated code
	"google.golang.org/grpc"
)

func main() {
	//env := os.Getenv("ENV")
	//if env == "" {
	//	log.Fatal("ENV environment variable not set")
	//}
	//err := godotenv.Load(env)
	err := config.NewConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lis, err := net.Listen("tcp", ":"+config.Cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	baseServer := grpc.NewServer()
	pb.RegisterDeployerServer(baseServer, NewServer())
	log.Printf("gRPC server started on port %s", config.Cfg.Port)
	baseServer.Serve(lis)
}
