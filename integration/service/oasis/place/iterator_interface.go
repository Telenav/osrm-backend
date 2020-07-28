package place

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
	"github.com/Telenav/osrm-backend/integration/util/osrmconnector"
)

// IteratorGenerator creates finders for different purpose, all finders must implement NearbyStationsIterator interface
type IteratorGenerator interface {
	// NewIterator4Orig creates finder to search for nearby charge stations near orig
	NewIterator4Orig(oasisReq *oasis.Request) iteratortype.NearbyStationsIterator

	// NewIterator4Dest creates finder to search for nearby charge stations near destination
	NewIterator4Dest(oasisReq *oasis.Request) iteratortype.NearbyStationsIterator

	// NewIterator4LowEnergyLocation creates finder to search for nearby charge stations when energy is low
	NewIterator4LowEnergyLocation(location *nav.Location) iteratortype.NearbyStationsIterator
}

// Algorithm contains algorithm implemented based on NearbyStationsIterator
type Algorithm interface {
	// FindOverlapBetweenStations finds overlap charge stations based on two iterator
	FindOverlapBetweenStations(iterF iteratortype.NearbyStationsIterator,
		iterS iteratortype.NearbyStationsIterator) []*iteratortype.ChargeStationInfo

	// CalcWeightBetweenChargeStationsPair accepts two iterators and calculates weights between each pair of iterators
	CalcWeightBetweenChargeStationsPair(from iteratortype.NearbyStationsIterator,
		to iteratortype.NearbyStationsIterator,
		table osrmconnector.TableRequster) ([]iteratortype.NeighborInfo, error)

	// CalculateWeightBetweenNeighbors accepts locations array, which will search for nearby
	// charge stations and then calculate weight between stations, the result is used to
	// construct graph.
	CalculateWeightBetweenNeighbors(locations []*nav.Location,
		oc *osrmconnector.OSRMConnector,
		finder IteratorGenerator) chan iteratortype.WeightBetweenNeighbors
}
