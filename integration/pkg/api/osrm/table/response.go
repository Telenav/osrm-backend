package table

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"

// Response represents OSRM api v1 table response.
type Response struct {
	Code         string            `json:"code"`
	Message      string            `json:"message,omitempty"`
	Sources      []*route.Waypoint `json:"sources"`
	Destinations []*route.Waypoint `json:"destinations"`
	Durations    [][]*float64      `json:"durations"`
	Distances    [][]*float64      `json:"distances"`
}
