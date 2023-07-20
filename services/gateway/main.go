package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	GeoServer "github.com/harlow/go-micro-services/services/gateway/src"
)

type server interface {
	Run(int) error
}

func main() {
	var (
		port        = flag.Int("port", 8080, "The service port")
		profileaddr = flag.String("profileaddr", "profile:8080", "Profile service addr")
		searchaddr  = flag.String("searchaddr", "search:8080", "Search service addr")
	)
	flag.Parse()

	var srv server

	srv = GeoServer.New(
		dial(*searchaddr),
		dial(*profileaddr),
	)

	if err := srv.Run(*port); err != nil {
		log.Fatalf("run src error: %v", err)
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
