package stationfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchhelper"
)

const lowEnergyLocationCandidateNumber = 20

type lowEnergyLocationStationFinder struct {
	location *nav.Location
	bf       *basicFinder
}

func NewLowEnergyLocationStationFinder(sc *searchconnector.TNSearchConnector, location *nav.Location) *lowEnergyLocationStationFinder {
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
