package haversine

import (
	"math"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371000
)

// GreatCircleDistance: Calculates the Haversine distance between two points in meters.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func GreatCircleDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1 = lat1 * (math.Pi / 180.0)
	lat2 = lat2 * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}
