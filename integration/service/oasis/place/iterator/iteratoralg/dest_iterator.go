package iteratoralg

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
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

func (di *destIterator) IterateNearbyStations() <-chan *iteratortype.ChargeStationInfo {
	c := make(chan *iteratortype.ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := iteratortype.ChargeStationInfo{
			ID:       iteratortype.DestLocationIDStr,
			Location: *di.location,
		}
		c <- &station
	}()

	return c
}
