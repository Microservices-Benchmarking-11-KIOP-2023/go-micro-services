package main

import (
	"flag"
	"fmt"
	"github.com/harlow/go-micro-services/services/search/src"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"google.golang.org/grpc"
)

type server interface {
	Run(int) error
}

func main() {
	var (
		port     = flag.Int("port", 8080, "The service port")
		geoaddr  = flag.String("geoaddr", "geo:8080", "Geo server addr")
		rateaddr = flag.String("rateaddr", "rate:8080", "Rate server addr")
	)
	flag.Parse()

	var srv server

	srv = src.New(
		dial(*geoaddr),
		dial(*rateaddr),
	)

	if err := srv.Run(*port); err != nil {
		log.Fatalf("run search error: %v", err)
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
