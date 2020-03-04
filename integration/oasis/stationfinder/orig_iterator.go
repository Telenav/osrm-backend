package stationfinder

// OrigLocationName defines name for orig
const OrigLocationName string = "orig_location"

type origIterator struct {
	location *StationCoordinate
}

func NewOrigIter(location *StationCoordinate) *origIterator {
	return &origIterator{
		location: location,
	}
}

func (oi *origIterator) iterateNearbyStations() <-chan ChargeStationInfo {
	c := make(chan ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := ChargeStationInfo{
			ID:       OrigLocationName,
			Location: *oi.location,
		}
		c <- station
	}()

	return c
}
