package stationhandler

type nearbyStationsIterator interface {
	iterateNearbyStations() <-chan ChargeStationInfo
}

// ChargeStationInfo defines charge station information
type ChargeStationInfo struct {
	id       string
	location StationCoordinate
	err      error
}

// Location returns coordinate of charge station
func (c ChargeStationInfo) Location() StationCoordinate {
	return c.location
}

// ID returns unique charge stations id
func (c ChargeStationInfo) ID() string {
	return c.id
}

// StationCoordinate represents location information
type StationCoordinate struct {
	Lat float64
	Lon float64
}
