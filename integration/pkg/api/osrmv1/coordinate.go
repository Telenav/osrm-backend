package osrmv1

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
)

// Coordinate represents lat/lon of a GPS point.
type Coordinate struct {
	Lat float64
	Lon float64
}

// Coordinates represents a list of GPS points.
type Coordinates []Coordinate

// String convert Coordinate to string. Lat/lon precision is 6.
func (c *Coordinate) String() string {

	s := fmt.Sprintf("%.6f%s%.6f", c.Lon, api.Comma, c.Lat)
	return s
}

// String convert Coordinates to string. Lat/lon precision is 6.
func (c *Coordinates) String() string {
	var s string
	for _, coord := range *c {
		if len(s) > 0 {
			s += api.Semicolon
		}
		s += coord.String()
	}

	return s
}
