package table

import (
	"github.com/Telenav/osrm-backend/integration/api/osrm/osrmtype"
)

// Response represents OSRM api v1 table response.
type Response struct {
	Code         string               `json:"code"`
	Message      string               `json:"message,omitempty"`
	Sources      []*osrmtype.Waypoint `json:"sources"`
	Destinations []*osrmtype.Waypoint `json:"destinations"`
	Durations    [][]float64          `json:"durations"`
	Distances    [][]float64          `json:"distances"`
}
