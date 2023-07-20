package src

import (
	"encoding/json"
	"fmt"
	profile "github.com/harlow/go-micro-services/services/profile/proto"
	"log"
	"net"

	"github.com/harlow/go-micro-services/services/profile/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// New returns a new server
func New() *Profile {
	return &Profile{
		profiles: loadProfiles("data/hotels.json"),
	}
}

// Profile implements the profile service
type Profile struct {
	profiles map[string]*profile.Hotel
}

// Run starts the server
func (s *Profile) Run(port int) error {
	srv := grpc.NewServer()
	profile.RegisterProfileServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return srv.Serve(lis)
}

// GetProfiles returns hotel profiles for requested IDs
func (s *Profile) GetProfiles(ctx context.Context, req *profile.Request) (*profile.Result, error) {
	res := new(profile.Result)
	for _, id := range req.HotelIds {
		res.Hotels = append(res.Hotels, s.getProfile(id))
	}
	return res, nil
}

func (s *Profile) getProfile(id string) *profile.Hotel {
	return s.profiles[id]
}

// loadProfiles loads hotel profiles from a JSON file.
func loadProfiles(path string) map[string]*profile.Hotel {
	var (
		file   = data.MustAsset(path)
		hotels []*profile.Hotel
	)

	if err := json.Unmarshal(file, &hotels); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}

	profiles := make(map[string]*profile.Hotel)
	for _, hotel := range hotels {
		profiles[hotel.Id] = hotel
	}
	return profiles
}
