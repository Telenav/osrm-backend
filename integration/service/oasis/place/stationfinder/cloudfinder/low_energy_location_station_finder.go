package cloudfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/searchhelper"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
)

// LowEnergyLocationCandidateNumber indicates how much charge station to be searched for low energy point
const LowEnergyLocationCandidateNumber = 20

type lowEnergyLocationStationFinder struct {
	location *nav.Location
	*basicFinder
}

func newLowEnergyLocationStationFinder(sc *searchconnector.TNSearchConnector, location *nav.Location) *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		location,
		newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (lFinder *lowEnergyLocationStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: lFinder.location.Lat,
			Lon: lFinder.location.Lon},
		LowEnergyLocationCandidateNumber,
		-1)
	lFinder.getNearbyChargeStations(req)
	return
}
