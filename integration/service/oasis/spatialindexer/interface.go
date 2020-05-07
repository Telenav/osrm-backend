package spatialindexer

import (
	"math"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
)

// Location for poi point
// todo codebear801 will be replaced by the one in nav
type Location struct {
	Lat float64
	Lon float64
}

// PointInfo records point related information such as ID and location
type PointInfo struct {
	ID       PointID
	Location Location
}

// RankedPointInfo used to record ranking result, distance to specific point could be used for ranking
type RankedPointInfo struct {
	PointInfo
	Distance float64
	Duration float64
}

// PointID defines ID for given point(location, point of interest)
// Only the data used for pre-processing contains valid PointID
type PointID int64

// String converts PointID to string
func (p PointID) String() string {
	return strconv.FormatInt((int64)(p), 10)
}

// UnlimitedCount means all spatial search result will be returned
const UnlimitedCount = math.MaxInt32

// Finder answers special query
type Finder interface {

	// FindNearByPointIDs returns a group of points near to given center location
	FindNearByPointIDs(center Location, radius float64, limitCount int) []*PointInfo
}

// Ranker used to ranking a group of points
type Ranker interface {

	// RankPointIDsByGreatCircleDistance ranks a group of points based on great circle distance to given location
	RankPointIDsByGreatCircleDistance(center Location, targets []*PointInfo) []*RankedPointInfo

	// RankPointIDsByShortestDistance ranks a group of points based on shortest path distance to given location
	RankPointIDsByShortestDistance(center Location, targets []*PointInfo) []*RankedPointInfo
}

// PlaceLocationQuerier returns *nav.location for given location
type PlaceLocationQuerier interface {

	// GetLocation returns *nav.Location for given placeID
	// Returns nil if given placeID is not found
	GetLocation(placeID string) *nav.Location
}

// PointsIterator provides iterateability for PointInfo
type PointsIterator interface {

	// IteratePoints returns channel for PointInfo
	IteratePoints() <-chan PointInfo
}
