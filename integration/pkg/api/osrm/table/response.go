package table

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"

// Response represents OSRM api v1 table response.
type Response struct {
	Code         string         `json:"code"`
	Message      string         `json:"message,omitempty"`
	Sources      []*Source      `json:"sources"`
	Destinations []*Destination `json:"destinations"`
	Durations    []*Duration    `json:"durations"`
	Distances    []*Distance    `json:"distances"`
}

// Source represents as way point object.  All sources will be listed in order.
type Source struct {
	route.Waypoint
}

// Destination represents as way point object.  All destination will be listed in order.
type Destination struct {
	route.Waypoint
}

// Duration gives the travel time from specific source to all other destinations
type Duration struct {
	Value []*float64
}

// Distance gives the travel distance from specific source to all other destinations
type Distance struct {
	Value []*float64
}
