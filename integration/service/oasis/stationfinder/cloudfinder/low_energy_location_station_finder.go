package cloudfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

// LowEnergyLocationCandidateNumber indicates how much charge station to be searched for low energy point
const LowEnergyLocationCandidateNumber = 20

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
		LowEnergyLocationCandidateNumber,
		-1)
	sf.bf.getNearbyChargeStations(req)
	return
}

// NearbyStationsIterator provides channel which contains near by station information for low energy location
func (sf *lowEnergyLocationStationFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	return sf.bf.IterateNearbyStations()
}
