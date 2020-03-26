package stationfinder

import "github.com/Telenav/osrm-backend/integration/pkg/api/nav"

// OrigLocationID defines name for orig
const OrigLocationID string = "orig_location"

type origIterator struct {
	location *nav.Location
}

func NewOrigIter(location *nav.Location) *origIterator {
	return &origIterator{
		location: location,
	}
}

func (oi *origIterator) iterateNearbyStations() <-chan ChargeStationInfo {
	c := make(chan ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := ChargeStationInfo{
			ID:       OrigLocationID,
			Location: *oi.location,
		}
		c <- station
	}()

	return c
}
