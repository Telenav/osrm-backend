package common

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
)

// PlaceInfo records place related information such as ID and location
type PlaceInfo struct {
	ID       PlaceID
	Location *nav.Location
}

// RankedPlaceInfo records PlaceInfo and information used for ranking
// e.g. use distance to specific place rank a group of PlaceInfo
type RankedPlaceInfo struct {
	PlaceInfo
	Weight *Weight
}

// PlaceID defines ID for given place(location, point of interest)
// The data used for pre-processing should contain valid PlaceID
type PlaceID int64

// String converts PlaceID to string
func (p PlaceID) String() string {
	return strconv.FormatInt((int64)(p), 10)
}
