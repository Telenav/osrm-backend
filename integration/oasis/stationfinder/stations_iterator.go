package stationfinder

type nearbyStationsIterator interface {
	iterateNearbyStations() <-chan ChargeStationInfo
}
