package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/harlow/go-micro-services/proto/auth"
	"github.com/harlow/go-micro-services/proto/geo"
	"github.com/harlow/go-micro-services/proto/profile"
	"github.com/harlow/go-micro-services/proto/rate"
	"github.com/harlow/go-micro-services/trace"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type inventory struct {
	Hotels    []*profile.Hotel `json:"hotels"`
	RatePlans []*rate.RatePlan `json:"ratePlans"`
}

// newServer returns a server with initialization data loaded.
func newServer(geoAddr, profileAddr, rateAddr *string) apiServer {
	return apiServer{
		serverName:    "api.v1",
		GeoClient:     geo.NewGeoClient(mustDial(geoAddr)),
		ProfileClient: profile.NewProfileClient(mustDial(profileAddr)),
		RateClient:    rate.NewRateClient(mustDial(rateAddr)),
	}
}

// mustDial ensures the tcp connection to specified address.
func mustDial(addr *string) *grpc.ClientConn {
	conn, err := grpc.Dial(*addr)
	if err != nil {
		log.Fatalf("dial failed: %s", err)
		panic(err)
	}
	return conn
}

type apiServer struct {
	serverName string
	geo.GeoClient
	profile.ProfileClient
	rate.RateClient
}

func (s apiServer) requestHandler(w http.ResponseWriter, r *http.Request) {
	t := trace.NewTracer()
	t.In("www", s.serverName)
	defer t.Out(s.serverName, "www", time.Now())

	// context and metadata
	md := metadata.Pairs("traceID", t.TraceID, "fromName", s.serverName)
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, md)

	// read and validate in/out arguments
	inDate := r.URL.Query().Get("inDate")
	outDate := r.URL.Query().Get("outDate")
	if inDate == "" || outDate == "" {
		http.Error(w, "Please specify inDate / outDate", http.StatusBadRequest)
		return
	}

	// get hotels within geo box
	reply, err := s.BoundedBox(ctx, &geo.Rectangle{
		Lo: &geo.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &geo.Point{Latitude: 420000000, Longitude: -730000000},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// make reqeusts for profiles and rates
	profileCh := s.getHotels(ctx, reply.HotelIds)
	rateCh := s.getRatePlans(ctx, reply.HotelIds, inDate, outDate)

	// wait on profiles reply
	profileReply := <-profileCh
	if err := profileReply.err; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// wait on rates reply
	rateReply := <-rateCh
	if err := rateReply.err; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// build the final inventory response
	inventory := inventory{
		Hotels:    profileReply.hotels,
		RatePlans: rateReply.ratePlans,
	}

	// encode JSON for rendering
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type rateResults struct {
	ratePlans []*rate.RatePlan
	err       error
}

func (s apiServer) getRatePlans(ctx context.Context, hotelIDs []int32, inDate string, outDate string) chan rateResults {
	ch := make(chan rateResults, 1)

	go func() {
		reply, err := s.GetRates(ctx,
			&rate.Args{
				HotelIds: hotelIDs,
				InDate:   inDate,
				OutDate:  outDate,
			})

		ch <- rateResults{
			ratePlans: reply.RatePlans,
			err:       err,
		}
	}()

	return ch
}

type profileResults struct {
	hotels []*profile.Hotel
	err    error
}

func (s apiServer) getHotels(ctx context.Context, hotelIDs []int32) chan profileResults {
	ch := make(chan profileResults, 1)

	go func() {
		reply, err := s.GetHotels(ctx, &profile.Args{HotelIds: hotelIDs})

		ch <- profileResults{
			hotels: reply.Hotels,
			err:    err,
		}
	}()

	return ch
}

func main() {
	var (
		port        = flag.String("port", "8080", "The server port")
		authAddr    = flag.String("auth", "auth:8080", "The Auth server address in the format of host:port")
		geoAddr     = flag.String("geo", "geo:8080", "The Geo server address in the format of host:port")
		profileAddr = flag.String("profile", "profile:8080", "The Pofile server address in the format of host:port")
		rateAddr    = flag.String("rate", "rate:8080", "The Rate Code server address in the format of host:port")
	)
	flag.Parse()

	server := newServer(geoAddr, profileAddr, rateAddr)
	authClient := auth.NewAuthClient(mustDial(authAddr))
	authHandler := NewAuthHandler(authClient)
	requesHandler := http.HandlerFunc(server.requestHandler)
	http.Handle("/", authHandler(requesHandler))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}