package place

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

// UnlimitedCount means all spatial search result will be returned
const UnlimitedCount = math.MaxInt32

// Finder answers spatial query
type Finder interface {

	// FindNearByPlaceIDs returns a group of places near to given center location
	FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*entity.PlaceInfo
}

// Ranker used to ranking a group of places
type Ranker interface {

	// RankPlaceIDsByGreatCircleDistance ranks a group of places based on great circle distance to given location
	RankPlaceIDsByGreatCircleDistance(center nav.Location, targets []*entity.PlaceInfo) []*entity.RankedPlaceInfo

	// RankPlaceIDsByShortestDistance ranks a group of places based on shortest path distance to given location
	RankPlaceIDsByShortestDistance(center nav.Location, targets []*entity.PlaceInfo) []*entity.RankedPlaceInfo
}

// LocationQuerier returns *nav.location for given location
type LocationQuerier interface {

	// GetLocation returns *nav.Location for given placeID
	// Returns nil if given placeID is not found
	GetLocation(placeID string) *nav.Location
}

// PlacesIterator provides iterateability for PlaceInfo
type PlacesIterator interface {

	// IteratePlaces returns channel for PlaceInfo
	IteratePlaces() <-chan entity.PlaceInfo
}
