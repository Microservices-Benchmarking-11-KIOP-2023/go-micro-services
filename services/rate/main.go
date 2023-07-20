package main

import (
	"flag"
	"fmt"
	"github.com/harlow/go-micro-services/rate/src"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type server interface {
	Run(int) error
}

func main() {
	var (
		port = flag.Int("port", 8080, "The service port")
	)
	flag.Parse()

	var srv server

	srv = src.New()

	if err := srv.Run(*port); err != nil {
		log.Fatalf("run rate error: %v", err)
	}
}

func dial(addr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}

	return conn
}
