package stationfindertype

import "github.com/Telenav/osrm-backend/integration/pkg/api/nav"

// NearbyStationsIterator provide interator for near by stations
type NearbyStationsIterator interface {
	// IterateNearbyStations returns a channel which contains near by charge station under certain conditions
	IterateNearbyStations() <-chan *ChargeStationInfo
}

// ChargeStationInfo defines charge station information
type ChargeStationInfo struct {
	ID       string
	Location nav.Location
	err      error
}

// Weight represent weight information
type Weight struct {
	Duration float64
	Distance float64
}

// NeighborInfo represent cost information between two charge stations
type NeighborInfo struct {
	FromID       string
	FromLocation nav.Location
	ToID         string
	ToLocation   nav.Location
	Weight
}

// WeightBetweenNeighbors contains a group of neighbors information
type WeightBetweenNeighbors struct {
	NeighborsInfo []NeighborInfo
	Err           error
}
