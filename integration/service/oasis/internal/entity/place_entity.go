package entity

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
)

// PlaceWithLocation records place related information such as ID and location
type PlaceWithLocation struct {
	ID       PlaceID
	Location *nav.Location
}

// TransferInfo records target PlaceWithLocation and Weight information used
// for ranking
// e.g. use distance to specific place rank a group of PlaceWithLocation
type TransferInfo struct {
	PlaceWithLocation
	Weight *Weight
}

// PlaceID defines ID for given place(location, point of interest)
// The data used for pre-processing must contain valid PlaceID, which means it
// either a int64 directly or be processed as int64
type PlaceID int64

// String converts PlaceID to string
func (p PlaceID) String() string {
	return strconv.FormatInt((int64)(p), 10)
}

// Weight represent weight information between stations
type Weight struct {
	Duration float64
	Distance float64
}
