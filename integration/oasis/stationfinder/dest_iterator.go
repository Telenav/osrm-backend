package stationfinder

import "github.com/Telenav/osrm-backend/integration/nav"

// DestLocationID defines name for dest
const DestLocationID string = "dest_location"

type destIterator struct {
	location *nav.Location
}

func NewDestIter(location *nav.Location) *destIterator {
	return &destIterator{
		location: location,
	}
}

func (di *destIterator) iterateNearbyStations() <-chan ChargeStationInfo {
	c := make(chan ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := ChargeStationInfo{
			ID:       DestLocationID,
			Location: *di.location,
		}
		c <- station
	}()

	return c
}
