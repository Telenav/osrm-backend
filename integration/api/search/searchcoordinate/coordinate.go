package searchcoordinate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
)

// Coordinate represents lat/lon of a GPS point.
type Coordinate struct {
	Lat float64
	Lon float64
}

// String convert Coordinate to string. Lat/lon precision is 6.
func (c *Coordinate) String() string {

	s := fmt.Sprintf("%.6f%s%.6f", c.Lat, api.Comma, c.Lon)
	return s
}

// ParseCoordinate parses string to coordinate.
func ParseCoordinate(str string) (Coordinate, error) {
	c := Coordinate{}

	splits := strings.Split(str, api.Comma)
	if len(splits) < 2 {
		return c, fmt.Errorf("parse %s to Coordinate failed", str)
	}

	var err error
	if c.Lat, err = strconv.ParseFloat(splits[0], 64); err != nil {
		return c, fmt.Errorf("parse Lon from %s failed", str)
	}
	if c.Lon, err = strconv.ParseFloat(splits[1], 64); err != nil {
		return c, fmt.Errorf("parse Lan from %s failed", str)
	}
	return c, nil
}
