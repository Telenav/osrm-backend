package spatialindexer

import (
	"math"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
)

// PlaceInfo records place related information such as ID and location
type PlaceInfo struct {
	ID       PlaceID
	Location nav.Location
}

// RankedPlaceInfo used to record ranking result, e.g. distance to specific place could be used for ranking
type RankedPlaceInfo struct {
	PlaceInfo
	Distance float64
	Duration float64
}

// PlaceID defines ID for given place(location, point of interest)
// Only the data used for pre-processing contains valid PlaceID
type PlaceID int64

// String converts PlaceID to string
func (p PlaceID) String() string {
	return strconv.FormatInt((int64)(p), 10)
}

// UnlimitedCount means all spatial search result will be returned
const UnlimitedCount = math.MaxInt32

// Finder answers special query
type Finder interface {

	// FindNearByPlaceIDs returns a group of places near to given center location
	FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*PlaceInfo
}

// Ranker used to ranking a group of places
type Ranker interface {

	// RankPlaceIDsByGreatCircleDistance ranks a group of places based on great circle distance to given location
	RankPlaceIDsByGreatCircleDistance(center nav.Location, targets []*PlaceInfo) []*RankedPlaceInfo

	// RankPlaceIDsByShortestDistance ranks a group of places based on shortest path distance to given location
	RankPlaceIDsByShortestDistance(center nav.Location, targets []*PlaceInfo) []*RankedPlaceInfo
}

// PlaceLocationQuerier returns *nav.location for given location
type PlaceLocationQuerier interface {

	// GetLocation returns *nav.Location for given placeID
	// Returns nil if given placeID is not found
	GetLocation(placeID string) *nav.Location
}

// PlacesIterator provides iterateability for PlaceInfo
type PlacesIterator interface {

	// IteratePlaces returns channel for PlaceInfo
	IteratePlaces() <-chan PlaceInfo
}
