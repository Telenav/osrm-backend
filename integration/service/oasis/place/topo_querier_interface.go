package place

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

// TopoQuerier used to return topological information of charge stations
type TopoQuerier interface {

	// GetConnectedPlaces finds connected places by given placeID and return them in recorded sequence
	// Returns nil if given placeID is not found or no connectivity
	GetConnectedPlaces(placeID entity.PlaceID) []*entity.TransferInfo

	// GetLocation returns location of given station id
	// Returns nil if given placeID is not found
	GetLocation(placeID entity.PlaceID) *nav.Location
}
