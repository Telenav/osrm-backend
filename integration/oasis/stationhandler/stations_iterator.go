package stationhandler

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
)

type nearbyStationsIterator interface {
	iterateNearbyStations() <-chan chargeStationInfo
}

type chargeStationInfo struct {
	id       string
	location coordinate.Coordinate
	err      error
}

func (c chargeStationInfo) Location() coordinate.Coordinate {
	return c.location
}

func buildChargeStationInfoDict(iter nearbyStationsIterator) map[string]bool {
	dict := make(map[string]bool)
	c := iter.iterateNearbyStations()
	for item := range c {
		dict[item.id] = true
	}

	return dict
}

func FindOverlapBetweenStations(iterF nearbyStationsIterator, iterS nearbyStationsIterator) []chargeStationInfo {
	overlap := make([]chargeStationInfo, 10)
	dict := buildChargeStationInfoDict(iterF)
	c := iterS.iterateNearbyStations()
	for item := range c {
		if _, has := dict[item.id]; has {
			overlap = append(overlap, item)
		}
	}

	return overlap
}
