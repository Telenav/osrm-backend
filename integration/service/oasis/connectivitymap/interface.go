package connectivitymap

import "github.com/Telenav/osrm-backend/integration/pkg/api/nav"

// QueryResult records topological query result
type QueryResult struct {
	StationID       string
	StationLocation *nav.Location
	Distance        float64
	Duration        float64
}

// Querier used to return topological information of charge stations
type Querier interface {

	// NearByStationQuery finds near by stations by given stationID and return them in recorded sequence
	// Returns nil if given stationID is not found or no connectivity
	NearByStationQuery(stationID string) []*QueryResult

	// GetLocation returns location of given station id
	// Returns nil if given stationID is not found
	GetLocation(stationID string) *nav.Location
}
