// Package coordinate defines OSRM request coordinates.
package coordinate

import (
	"fmt"
	"strconv"
	"strings"

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

// ParseCoordinate parses string to coorinate.
func ParseCoordinate(str string) (Coordinate, error) {
	c := Coordinate{}

	splits := strings.Split(str, api.Comma)
	if len(splits) < 2 {
		return c, fmt.Errorf("parse %s to Coordinate failed", str)
	}

	var err error
	if c.Lon, err = strconv.ParseFloat(splits[0], 64); err != nil {
		return c, fmt.Errorf("parse Lon from %s failed", str)
	}
	if c.Lat, err = strconv.ParseFloat(splits[1], 64); err != nil {
		return c, fmt.Errorf("parse Lan from %s failed", str)
	}
	return c, nil
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

// ParseCoordinates parses string to coordinates.
func ParseCoordinates(str string) (Coordinates, error) {
	var coordinates Coordinates
	splits := strings.Split(str, api.Semicolon)
	for _, s := range splits {
		c, err := ParseCoordinate(s)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, c)
	}
	return coordinates, nil
}
