package main

import (
	"flag"
	"fmt"
	FrontendServer "github.com/harlow/go-micro-services/services/frontend"
	GeoServer "github.com/harlow/go-micro-services/services/geo"
	ProfileServer "github.com/harlow/go-micro-services/services/profile"
	RateServer "github.com/harlow/go-micro-services/services/rate"
	SearchServer "github.com/harlow/go-micro-services/services/search"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"

	"google.golang.org/grpc"
)

type server interface {
	Run(int) error
}

func main() {
	var (
		port        = flag.Int("port", 8080, "The service port")
		profileaddr = flag.String("profileaddr", "profile:8080", "Profile service addr")
		geoaddr     = flag.String("geoaddr", "geo:8080", "Geo server addr")
		rateaddr    = flag.String("rateaddr", "rate:8080", "Rate server addr")
		searchaddr  = flag.String("searchaddr", "search:8080", "Search service addr")
	)
	flag.Parse()

	var srv server
	var cmd = os.Args[1]

	switch cmd {
	case "geo":
		srv = GeoServer.New()
	case "rate":
		srv = RateServer.New()
	case "profile":
		srv = ProfileServer.New()
	case "search":
		srv = SearchServer.New(
			dial(*geoaddr),
			dial(*rateaddr),
		)
	case "frontend":
		srv = FrontendServer.New(
			dial(*searchaddr),
			dial(*profileaddr),
		)
	default:
		log.Fatalf("unknown cmd: %s", cmd)
	}

	if err := srv.Run(*port); err != nil {
		log.Fatalf("run %s error: %v", cmd, err)
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
