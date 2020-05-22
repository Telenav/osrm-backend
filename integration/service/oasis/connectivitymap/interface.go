package connectivitymap

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

// Querier used to return topological information of charge stations
type Querier interface {

	// NearByStationQuery finds near by stations by given placeID and return them in recorded sequence
	// Returns nil if given placeID is not found or no connectivity
	NearByStationQuery(placeID common.PlaceID) []*common.RankedPlaceInfo

	// GetLocation returns location of given station id
	// Returns nil if given placeID is not found
	GetLocation(placeID common.PlaceID) *nav.Location
}
