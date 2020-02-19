package stationfinder

type nearbyStationsIterator interface {
	iterateNearbyStations() <-chan ChargeStationInfo
}

// ChargeStationInfo defines charge station information
type ChargeStationInfo struct {
	ID       string
	Location StationCoordinate
	err      error
}

// StationCoordinate represents location information
type StationCoordinate struct {
	Lat float64
	Lon float64
}
