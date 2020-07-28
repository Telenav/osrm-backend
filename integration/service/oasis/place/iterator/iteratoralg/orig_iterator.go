package iteratoralg

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
)

type origIterator struct {
	location *nav.Location
}

// NewDestIter creates origIterator
// origIterator wraps single point of orig which adopts algorithms' requirement
func NewOrigIter(location *nav.Location) *origIterator {
	return &origIterator{
		location: location,
	}
}

func (oi *origIterator) IterateNearbyStations() <-chan *iteratortype.ChargeStationInfo {
	c := make(chan *iteratortype.ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := iteratortype.ChargeStationInfo{
			ID:       iteratortype.OrigLocationIDStr,
			Location: *oi.location,
		}
		c <- &station
	}()

	return c
}
