package table

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"

// Response represents OSRM api v1 table response.
type Response struct {
	Code         string         `json:"code"`
	Message      string         `json:"message,omitempty"`
	Sources      []*Source      `json:"sources"`
	Destinations []*Destination `json:"destinations"`
	Durations    [][]*float64   `json:"durations"`
	Distances    [][]*float64   `json:"distances"`
}

// Source represents as way point object.  All sources will be listed in order.
type Source struct {
	route.Waypoint
}

// Destination represents as way point object.  All destination will be listed in order.
type Destination struct {
	route.Waypoint
}
