package stationfinder

import (
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
)

const lowEnergyLocationCandidateNumber = 20

type lowEnergyLocationStationFinder struct {
	location *StationCoordinate
	bf       *basicFinder
}

func NewLowEnergyLocationStationFinder(sc *searchconnector.TNSearchConnector, location *StationCoordinate) *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		location: location,
		bf:       newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (sf *lowEnergyLocationStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: sf.location.Lat,
			Lon: sf.location.Lon},
		lowEnergyLocationCandidateNumber,
		-1)
	sf.bf.getNearbyChargeStations(req)
	return
}

func (sf *lowEnergyLocationStationFinder) iterateNearbyStations() <-chan ChargeStationInfo {
	return sf.bf.iterateNearbyStations()
}
