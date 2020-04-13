package stationfinderalg

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

type destIterator struct {
	location *nav.Location
}

// NewDestIter creates destIterator
// destIterator wraps single point of dest which adopts algorithms' requirement
func NewDestIter(location *nav.Location) *destIterator {
	return &destIterator{
		location: location,
	}
}

func (di *destIterator) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	c := make(chan *stationfindertype.ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := stationfindertype.ChargeStationInfo{
			ID:       stationfindertype.DestLocationID,
			Location: *di.location,
		}
		c <- &station
	}()

	return c
}
