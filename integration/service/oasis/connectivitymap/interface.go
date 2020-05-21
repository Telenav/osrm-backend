package connectivitymap

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

// Querier used to return topological information of charge stations
type Querier interface {

	// NearByStationQuery finds near by stations by given stationID and return them in recorded sequence
	// Returns nil if given stationID is not found or no connectivity
	NearByStationQuery(stationID string) []*common.RankedPlaceInfo

	// GetLocation returns location of given station id
	// Returns nil if given stationID is not found
	GetLocation(stationID string) *nav.Location
}
