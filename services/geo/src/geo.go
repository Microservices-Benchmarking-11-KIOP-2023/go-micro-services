package src

import (
	"encoding/json"
	"fmt"
	geo "github.com/harlow/go-micro-services/geo/proto"
	"log"
	"math"
	"net"

	"github.com/harlow/go-micro-services/geo/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	maxSearchRadius = 10
	earthRadiusKm   = 6371
)

// point represents a hotel's geographic location on map
type point struct {
	Pid  string  `json:"hotelId"`
	Plat float64 `json:"lat"`
	Plon float64 `json:"lon"`
}

// Implement Point interface
func (p *point) Lat() float64 { return p.Plat }
func (p *point) Lon() float64 { return p.Plon }
func (p *point) Id() string   { return p.Pid }

// New returns a new server
func New() *Geo {
	return &Geo{
		points: loadPoints("data/geo.json"),
	}
}

// Server implements the geo service
type Geo struct {
	points []*point
}

// Run starts the server
func (s *Geo) Run(port int) error {
	srv := grpc.NewServer()
	geo.RegisterGeoServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return srv.Serve(lis)
}

// Haversine formula to calculate the distance between two points on the earth
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1Rad, lon1Rad := lat1*math.Pi/180, lon1*math.Pi/180
	lat2Rad, lon2Rad := lat2*math.Pi/180, lon2*math.Pi/180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// Nearby returns all hotels within a given distance.
func (s *Geo) Nearby(ctx context.Context, req *geo.Request) (*geo.Result, error) {
	var nearbyHotels []string
	for _, point := range s.points {
		distance := haversineDistance(float64(req.Lat), float64(req.Lon), point.Plat, point.Plon)
		if distance <= maxSearchRadius {
			nearbyHotels = append(nearbyHotels, point.Pid)
		}
	}

	return &geo.Result{HotelIds: nearbyHotels}, nil
}

func loadPoints(path string) []*point {
	var (
		file   = data.MustAsset(path)
		points []*point
	)

	if err := json.Unmarshal(file, &points); err != nil {
		log.Fatalf("Failed to load hotels: %v", err)
	}

	return points
}
