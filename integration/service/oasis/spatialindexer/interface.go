package spatialindexer

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

// UnlimitedCount means all spatial search result will be returned
const UnlimitedCount = math.MaxInt32

// Finder answers spatial query
type Finder interface {

	// FindNearByPlaceIDs returns a group of places near to given center location
	FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*common.PlaceInfo
}

// Ranker used to ranking a group of places
type Ranker interface {

	// RankPlaceIDsByGreatCircleDistance ranks a group of places based on great circle distance to given location
	RankPlaceIDsByGreatCircleDistance(center nav.Location, targets []*common.PlaceInfo) []*common.RankedPlaceInfo

	// RankPlaceIDsByShortestDistance ranks a group of places based on shortest path distance to given location
	RankPlaceIDsByShortestDistance(center nav.Location, targets []*common.PlaceInfo) []*common.RankedPlaceInfo
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
	IteratePlaces() <-chan common.PlaceInfo
}
