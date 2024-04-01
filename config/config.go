package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	TargetAddr string
	ClientName string
}

var Cfg *Config
var (
	DefaultClient = "localhost"
)

func NewConfig(filename ...string) error {
	if len(filename) == 0 {
		filename = append(filename, ".env")
	}

	if len(filename) > 1 {
		return fmt.Errorf("cannot pass more than 1 filename")
	}

	if err := godotenv.Load(filename...); err != nil {
		return fmt.Errorf("no .env file found")
	}

	targetAddr := os.Getenv("TARGET_CONTAINER_ADDR")
	if targetAddr == "" {
		return fmt.Errorf("TARGET_CONTAINER_ADDR environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("PORT environment variable not set")
	}


	clientName := os.Getenv("CLIENT_NAME")
	if clientName  == "" {
		clientName = DefaultClient
	}

	Cfg = &Config{
		TargetAddr:         targetAddr,
		Port: port,
		ClientName: clientName,
	} 

	return nil
}
