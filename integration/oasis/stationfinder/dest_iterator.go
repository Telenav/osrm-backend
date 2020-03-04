package stationfinder

const destLocationID string = "dest_location"

type destIterator struct {
	location *StationCoordinate
}

func NewDestIter(location *StationCoordinate) *destIterator {
	return &destIterator{
		location: location,
	}
}

func (di *destIterator) iterateNearbyStations() <-chan ChargeStationInfo {
	c := make(chan ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := ChargeStationInfo{
			ID:       destLocationID,
			Location: *di.location,
		}
		c <- station
	}()

	return c
}
