package frontend

import (
	"encoding/json"
	"fmt"
	profile "github.com/harlow/go-micro-services/services/profile/proto"
	search "github.com/harlow/go-micro-services/services/search/proto"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

// New returns a new server
func New(searchconn, profileconn *grpc.ClientConn) *Frontend {
	return &Frontend{
		searchClient:  search.NewSearchClient(searchconn),
		profileClient: profile.NewProfileClient(profileconn),
	}
}

// Frontend implements frontend service
type Frontend struct {
	searchClient  search.SearchClient
	profileClient profile.ProfileClient
}

// Run the server
func (s *Frontend) Run(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/hotels", http.HandlerFunc(s.searchHandler))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// Run the server
func (s *Frontend) searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	// in/out dates from query params
	inDate, outDate := r.URL.Query().Get("inDate"), r.URL.Query().Get("outDate")
	if inDate == "" || outDate == "" {
		http.Error(w, "Please specify inDate/outDate params", http.StatusBadRequest)
		return
	}

	// get lat/lon from query params
	latParam, lonParam := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	if latParam == "" || lonParam == "" {
		http.Error(w, "Please specify lat/lon params", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(latParam), 32)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(strings.TrimSpace(lonParam), 32)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	// search for best hotels
	searchResp, err := s.searchClient.Nearby(ctx, &search.NearbyRequest{
		Lat:     float32(lat),
		Lon:     float32(lon),
		InDate:  inDate,
		OutDate: outDate,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// grab locale from query params or default to en
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}

	// hotel profiles
	profileResp, err := s.profileClient.GetProfiles(ctx, &profile.Request{
		HotelIds: searchResp.HotelIds,
		Locale:   locale,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(geoJSONResponse(profileResp.Hotels))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// return a geoJSON response that allows google map to plot points directly on map
func geoJSONResponse(hs []*profile.Hotel) map[string]interface{} {
	var fs []interface{}

	for _, h := range hs {
		fs = append(fs, map[string]interface{}{
			"type": "Feature",
			"id":   h.Id,
			"properties": map[string]string{
				"name":         h.Name,
				"phone_number": h.PhoneNumber,
			},
			"geometry": map[string]interface{}{
				"type": "Point",
				"coordinates": []float32{
					h.Address.Lon,
					h.Address.Lat,
				},
			},
		})
	}

	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": fs,
	}
}