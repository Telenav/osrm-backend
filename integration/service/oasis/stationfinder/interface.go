package stationfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

// StationFinder creates finders for different purpose, all finders must implement NearbyStationsIterator interface
type StationFinder interface {
	// NewOrigStationFinder creates finder to search for nearby charge stations near orig
	NewOrigStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator

	// NewDestStationFinder creates finder to search for nearby charge stations near destination
	NewDestStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator

	// NewLowEnergyLocationStationFinder creates finder to search for nearby charge stations when energy is low
	NewLowEnergyLocationStationFinder(location *nav.Location) stationfindertype.NearbyStationsIterator
}

// Algorithm contains algorithm implemented based on NearbyStationsIterator
type Algorithm interface {
	// FindOverlapBetweenStations finds overlap charge stations based on two iterator
	FindOverlapBetweenStations(iterF stationfindertype.NearbyStationsIterator,
		iterS stationfindertype.NearbyStationsIterator) []*stationfindertype.ChargeStationInfo

	// CalcWeightBetweenChargeStationsPair accepts two iterators and calculates weights between each pair of iterators
	CalcWeightBetweenChargeStationsPair(from stationfindertype.NearbyStationsIterator,
		to stationfindertype.NearbyStationsIterator,
		table osrmconnector.TableRequster) ([]stationfindertype.NeighborInfo, error)

	// CalculateWeightBetweenNeighbors accepts locations array, which will search for nearby
	// charge stations and then calculate weight between stations, the result is used to
	// construct graph.
	CalculateWeightBetweenNeighbors(locations []*nav.Location,
		oc *osrmconnector.OSRMConnector,
		finder StationFinder) chan stationfindertype.WeightBetweenNeighbors
}
